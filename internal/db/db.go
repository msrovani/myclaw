package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Config holds SQLite configuration.
type Config struct {
	BaseDataDir  string // Root dir for multi-tenant databases, defaults to "data"
	Path         string // Overridden by DB Manager dynamically
	BusyTimeout  int    // milliseconds
	WALEnabled   bool
	VecEnabled   bool
	VecDimension int
	MaxReaders   int
	Env          string // dev, prod, edge
}

// DB provides concurrency-safe SQLite access with serialized writer
// and pooled readers. WAL mode is always preferred.
type DB struct {
	writer    *sql.DB
	readers   *sql.DB
	writeMu   sync.Mutex
	writeCh   chan writeOp
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	cfg       Config
	vecLoaded bool
}

type writeOp struct {
	fn     func(tx *sql.Tx) error
	result chan error
}

// Open creates a new DB instance with WAL, serialized writer, and reader pool.
func Open(cfg Config) (*DB, error) {
	// Writer connection — single, serialized.
	writerDSN := fmt.Sprintf("file:%s?_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(1)",
		cfg.Path, cfg.BusyTimeout)

	writer, err := sql.Open("sqlite", writerDSN)
	if err != nil {
		return nil, fmt.Errorf("db: open writer: %w", err)
	}
	writer.SetMaxOpenConns(1) // Serialized writer
	writer.SetMaxIdleConns(1)

	if err := writer.Ping(); err != nil {
		writer.Close()
		return nil, fmt.Errorf("db: ping writer: %w", err)
	}

	var mode string
	if err := writer.QueryRow("PRAGMA journal_mode=WAL").Scan(&mode); err != nil || mode != "wal" {
		writer.Close()
		return nil, fmt.Errorf("db: enable WAL failed: %w (got %q)", err, mode)
	}

	// Reader pool — multiple concurrent readers.
	readerDSN := fmt.Sprintf("file:%s?mode=ro&_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		cfg.Path, cfg.BusyTimeout)

	maxReaders := cfg.MaxReaders
	if maxReaders <= 0 {
		maxReaders = 4
	}

	readers, err := sql.Open("sqlite", readerDSN)
	if err != nil {
		writer.Close()
		return nil, fmt.Errorf("db: open readers: %w", err)
	}
	readers.SetMaxOpenConns(maxReaders)
	readers.SetMaxIdleConns(maxReaders)

	ctx, cancel := context.WithCancel(context.Background())
	d := &DB{
		writer:  writer,
		readers: readers,
		writeCh: make(chan writeOp, 256),
		ctx:     ctx,
		cancel:  cancel,
		cfg:     cfg,
	}

	// Start the serialized writer goroutine.
	d.wg.Add(1)
	go d.writerLoop()

	slog.Info("db: opened",
		"path", cfg.Path,
		"wal", cfg.WALEnabled,
		"max_readers", maxReaders,
		"busy_timeout", cfg.BusyTimeout,
	)

	return d, nil
}

// writerLoop processes write operations sequentially.
func (d *DB) writerLoop() {
	defer d.wg.Done()
	for {
		select {
		case <-d.ctx.Done():
			// Drain remaining operations.
			for op := range d.writeCh {
				op.result <- fmt.Errorf("db: shutting down")
			}
			return
		case op, ok := <-d.writeCh:
			if !ok {
				return
			}
			tx, err := d.writer.BeginTx(d.ctx, nil)
			if err != nil {
				op.result <- fmt.Errorf("db: begin tx: %w", err)
				continue
			}
			if err := op.fn(tx); err != nil {
				tx.Rollback()
				op.result <- err
			} else {
				op.result <- tx.Commit()
			}
		}
	}
}

// Write submits a write operation to the serialized writer queue.
// The function runs inside a transaction. Blocks until complete.
func (d *DB) Write(ctx context.Context, fn func(tx *sql.Tx) error) error {
	op := writeOp{
		fn:     fn,
		result: make(chan error, 1),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-d.ctx.Done():
		return fmt.Errorf("db: closed")
	case d.writeCh <- op:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-op.result:
		return err
	}
}

// Read executes a read-only query using the reader pool.
func (d *DB) Read(ctx context.Context, fn func(db *sql.DB) error) error {
	return fn(d.readers)
}

// ReadRow is a convenience for reading a single row.
func (d *DB) ReadRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.readers.QueryRowContext(ctx, query, args...)
}

// ReadRows is a convenience for reading multiple rows.
func (d *DB) ReadRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.readers.QueryContext(ctx, query, args...)
}

// Writer returns the raw writer DB for migrations. Use with caution.
func (d *DB) Writer() *sql.DB {
	return d.writer
}

// Checkpoint forces a WAL checkpoint.
func (d *DB) Checkpoint(ctx context.Context) error {
	_, err := d.writer.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		return fmt.Errorf("db: checkpoint: %w", err)
	}
	slog.Info("db: WAL checkpoint completed")
	return nil
}

// Stats returns database pool statistics.
type Stats struct {
	Writer        sql.DBStats `json:"writer"`
	Readers       sql.DBStats `json:"readers"`
	WriteQueueLen int         `json:"write_queue_len"`
	WriteQueueCap int         `json:"write_queue_cap"`
}

func (d *DB) Stats() Stats {
	return Stats{
		Writer:        d.writer.Stats(),
		Readers:       d.readers.Stats(),
		WriteQueueLen: len(d.writeCh),
		WriteQueueCap: cap(d.writeCh),
	}
}

// Close shuts down the database.
func (d *DB) Close() error {
	d.cancel()
	close(d.writeCh)
	d.wg.Wait()

	// Final checkpoint before close.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	d.Checkpoint(ctx)

	var errs []error
	if err := d.writer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("writer close: %w", err))
	}
	if err := d.readers.Close(); err != nil {
		errs = append(errs, fmt.Errorf("readers close: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("db: close errors: %v", errs)
	}
	slog.Info("db: closed")
	return nil
}
