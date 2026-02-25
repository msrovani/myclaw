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
	writeCh   chan writeOp
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	cfg       Config
	mu        sync.RWMutex
	closed    bool
}

type Stats struct {
	WriteQueueLen int
	WriteQueueCap int
}

type writeOp struct {
	fn     func(tx *sql.Tx) error
	result chan error
}

// Open creates a new DB instance with WAL, serialized writer, and reader pool.
func Open(cfg Config) (*DB, error) {
	// Writer connection — single, serialized.
	writerDSN := fmt.Sprintf("file:%s?_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(1)&_pragma=cache_size(-20000)&_pragma=mmap_size(268435456)&_pragma=temp_store(memory)",
		cfg.Path, cfg.BusyTimeout)

	writer, err := sql.Open("sqlite", writerDSN)
	if err != nil {
		return nil, fmt.Errorf("db: open writer: %w", err)
	}
	writer.SetMaxOpenConns(1)
	writer.SetMaxIdleConns(1)

	if err := writer.Ping(); err != nil {
		writer.Close()
		return nil, fmt.Errorf("db: ping writer: %w", err)
	}

	// Reader pool — multiple concurrent readers.
	readerDSN := fmt.Sprintf("file:%s?mode=ro&_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-20000)&_pragma=mmap_size(268435456)&_pragma=temp_store(memory)",
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

	d.wg.Add(1)
	go d.writerLoop()

	slog.Info("db: opened",
		"path", cfg.Path,
		"max_readers", maxReaders,
		"busy_timeout", cfg.BusyTimeout,
	)

	return d, nil
}

func (d *DB) writerLoop() {
	defer d.wg.Done()
	for {
		select {
		case <-d.ctx.Done():
			// Drain remaining operations.
			// Note: We don't close writeCh to avoid panics in Write().
			// We just stop processing once the context is cancelled.
			return
		case op := <-d.writeCh:
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
func (d *DB) Write(ctx context.Context, fn func(tx *sql.Tx) error) error {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return fmt.Errorf("db: closed")
	}
	d.mu.RUnlock()

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

func (d *DB) ReadRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.readers.QueryRowContext(ctx, query, args...)
}

func (d *DB) ReadRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.readers.QueryContext(ctx, query, args...)
}

func (d *DB) Writer() *sql.DB {
	return d.writer
}

func (d *DB) Checkpoint(ctx context.Context) error {
	_, err := d.writer.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		return fmt.Errorf("db: checkpoint: %w", err)
	}
	return nil
}

func (d *DB) Stats() Stats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return Stats{
		WriteQueueLen: len(d.writeCh),
		WriteQueueCap: cap(d.writeCh),
	}
}

func (d *DB) Close() error {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return nil
	}
	d.closed = true
	d.mu.Unlock()

	d.cancel()
	// We don't close(d.writeCh) here to avoid racing with Write()
	d.wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = d.Checkpoint(ctx)

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
