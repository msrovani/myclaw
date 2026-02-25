package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
)

func TestManager_CrossTenantIsolation(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		BaseDataDir: dir,
		BusyTimeout: 1000,
		MaxReaders:  1,
	}
	m := NewManager(cfg)
	defer m.CloseAll()

	ctxA := core.WithTenant(context.Background(), core.TenantContext{
		UID:         "user_a",
		WorkspaceID: "ws_1",
	})
	ctxB := core.WithTenant(context.Background(), core.TenantContext{
		UID:         "user_b",
		WorkspaceID: "ws_1", // Same workspace ID, different user
	})

	// User A DB
	dbA, err := m.GetDB(ctxA)
	if err != nil {
		t.Fatalf("GetDB user A: %v", err)
	}

	// Insert data into User A's BD
	err = dbA.Write(ctxA, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO memories (id, uid, workspace_id, content) VALUES ('m1', 'user_a', 'ws_1', 'secret')")
		return err
	})
	if err != nil {
		t.Fatalf("user A write: %v", err)
	}

	// User B DB
	dbB, err := m.GetDB(ctxB)
	if err != nil {
		t.Fatalf("GetDB user B: %v", err)
	}

	// Attempt to read User A's secret from User B's DB
	var content string
	err = dbB.ReadRow(ctxB, "SELECT content FROM memories WHERE id = 'm1'").Scan(&content)
	if err != sql.ErrNoRows {
		t.Fatalf("cross-tenant leakage! expected ErrNoRows, got err=%v content=%v", err, content)
	}

	// Verify paths are completely different
	expectedPathA := filepath.Join(dir, "tenants", "user_a", "workspaces", "ws_1", "memory.db")
	expectedPathB := filepath.Join(dir, "tenants", "user_b", "workspaces", "ws_1", "memory.db")

	if _, err := os.Stat(expectedPathA); err != nil {
		t.Errorf("user A db path missing: %v", err)
	}
	if _, err := os.Stat(expectedPathB); err != nil {
		t.Errorf("user B db path missing: %v", err)
	}
}

func TestManager_MissingContext(t *testing.T) {
	m := NewManager(Config{BaseDataDir: t.TempDir()})
	defer m.CloseAll()

	_, err := m.GetDB(context.Background())
	if err == nil {
		t.Fatal("GetDB should fail without TenantContext")
	}
}

func TestManager_PathTraversal(t *testing.T) {
	m := NewManager(Config{BaseDataDir: t.TempDir()})
	defer m.CloseAll()

	ctxHack := core.WithTenant(context.Background(), core.TenantContext{
		UID:         "../../etc",
		WorkspaceID: "passwd",
	})

	if _, err := m.GetDB(ctxHack); err == nil {
		t.Error("expected error on path traversal attempt")
	}
}
