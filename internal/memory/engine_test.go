package memory

import (
	"context"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
)

// Minimal test to explicitly save tokens while ensuring the engine compiles and works.
func TestEngine_Minimal(t *testing.T) {
	mgr := db.NewManager(db.Config{BaseDataDir: t.TempDir()})
	defer mgr.CloseAll()

	// Using nil provider for minimal unit test (engine handles nil gracefully)
	engine := NewEngine(mgr, nil)
	ctx := core.WithTenant(context.Background(), core.TenantContext{UID: "user1", WorkspaceID: "w1"})

	// 1. Add
	mem, err := engine.AddMemory(ctx, "Minimal token-saving test", AddOptions{Layer: LayerShortTerm})
	if err != nil {
		t.Fatalf("AddMemory: %v", err)
	}

	if mem.Layer != LayerShortTerm || mem.Content != "Minimal token-saving test" {
		t.Errorf("AddMemory mismatch")
	}

	// 2. Get
	fetched, err := engine.GetMemory(ctx, mem.ID)
	if err != nil {
		t.Fatalf("GetMemory: %v", err)
	}
	if fetched.ID != mem.ID {
		t.Errorf("GetMemory ID mismatch")
	}

	// 3. Update
	updated, err := engine.UpdateMemory(ctx, mem.ID, "Updated content", map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("UpdateMemory: %v", err)
	}
	if updated.Content != "Updated content" {
		t.Errorf("Updated content mismatch")
	}

	// 4. Delete
	if err := engine.DeleteMemory(ctx, mem.ID); err != nil {
		t.Fatalf("DeleteMemory: %v", err)
	}
}
