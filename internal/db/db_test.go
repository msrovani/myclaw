package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempDBPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.db")
}

func TestOpen_And_Close(t *testing.T) {
	path := tempDBPath(t)
	db, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		WALEnabled:  true,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	// Verify WAL mode
	var journalMode string
	db.ReadRow(context.Background(), "PRAGMA journal_mode").Scan(&journalMode)
	if journalMode != "wal" {
		t.Errorf("journal_mode = %q, want %q", journalMode, "wal")
	}
}

func TestDB_Write_Read(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		WALEnabled:  true,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	// Create table and insert via serialized writer.
	err = d.Write(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, val TEXT)")
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO test (id, val) VALUES (1, 'hello')")
		return err
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	// Read via reader pool.
	var val string
	err = d.ReadRow(context.Background(), "SELECT val FROM test WHERE id = ?", 1).Scan(&val)
	if err != nil {
		t.Fatalf("ReadRow: %v", err)
	}
	if val != "hello" {
		t.Errorf("val = %q, want %q", val, "hello")
	}
}

func TestDB_Write_ContextCancelled(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	err = d.Write(ctx, func(tx *sql.Tx) error {
		return nil
	})
	if err == nil {
		t.Error("Write should fail with cancelled context")
	}
}

func TestDB_Checkpoint(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	err = d.Checkpoint(context.Background())
	if err != nil {
		t.Errorf("Checkpoint: %v", err)
	}
}

func TestDB_Stats(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	stats := d.Stats()
	if stats.WriteQueueCap != 256 {
		t.Errorf("WriteQueueCap = %d, want 256", stats.WriteQueueCap)
	}
}

func TestMigrations(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	// Run core migrations.
	err = Migrate(d.Writer(), CoreMigrations())
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	// Verify tables exist.
	tables := []string{"memories", "memories_fts", "token_usage", "schema_migrations"}
	for _, table := range tables {
		var name string
		err := d.ReadRow(context.Background(), "SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}

	// Run again — should be idempotent.
	err = Migrate(d.Writer(), CoreMigrations())
	if err != nil {
		t.Fatalf("Migrate (idempotent): %v", err)
	}

	// Verify version recorded.
	var version int
	d.ReadRow(context.Background(), "SELECT MAX(version) FROM schema_migrations").Scan(&version)
	if version != 1 {
		t.Errorf("version = %d, want 1", version)
	}
}

func TestMigrations_FailedMigration(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	badMigrations := []Migration{
		{Version: 1, Description: "bad", Up: "INVALID SQL SYNTAX HERE !!!"},
	}

	err = Migrate(d.Writer(), badMigrations)
	if err == nil {
		t.Fatal("Migrate should fail for invalid SQL")
	}
}

func TestDB_ConcurrentReads(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  4,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	// Create test data.
	d.Write(context.Background(), func(tx *sql.Tx) error {
		tx.Exec("CREATE TABLE nums (n INTEGER)")
		for i := 0; i < 100; i++ {
			tx.Exec("INSERT INTO nums VALUES (?)", i)
		}
		return nil
	})

	// Concurrent reads — should not error.
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			var count int
			err := d.ReadRow(context.Background(), "SELECT COUNT(*) FROM nums").Scan(&count)
			if err != nil {
				done <- err
				return
			}
			if count != 100 {
				done <- err
			}
			done <- nil
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("concurrent read %d: %v", i, err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("concurrent reads timed out")
		}
	}
}

func TestDB_FileCreated(t *testing.T) {
	path := tempDBPath(t)
	d, err := Open(Config{
		Path:        path,
		BusyTimeout: 5000,
		MaxReaders:  2,
		Env:         "dev",
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	d.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("database file was not created")
	}
}
