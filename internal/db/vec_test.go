package db

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
)

func TestFloat32BytesConversion(t *testing.T) {
	vec := []float32{0.1, -0.2, 0.334, 1.0}
	b, err := Float32ToBytes(vec)
	if err != nil {
		t.Fatalf("Float32ToBytes failed: %v", err)
	}

	vec2, err := BytesToFloat32(b)
	if err != nil {
		t.Fatalf("BytesToFloat32 failed: %v", err)
	}

	if len(vec) != len(vec2) {
		t.Fatalf("length mismatch: %d vs %d", len(vec), len(vec2))
	}

	for i := range vec {
		if vec[i] != vec2[i] {
			t.Errorf("mismatch at index %d: %f != %f", i, vec[i], vec2[i])
		}
	}
}

func TestCosineDistance(t *testing.T) {
	a := []float32{1, 0, 0}
	b := []float32{1, 0, 0}
	c := []float32{0, 1, 0}
	d := []float32{-1, 0, 0}

	if d1 := CosineDistance(a, b); math.Abs(float64(d1-0.0)) > 1e-6 {
		t.Errorf("expected 0 for identical vectors, got %f", d1)
	}
	if d2 := CosineDistance(a, c); math.Abs(float64(d2-1.0)) > 1e-6 {
		t.Errorf("expected 1.0 for orthogonal vectors, got %f", d2)
	}
	if d3 := CosineDistance(a, d); math.Abs(float64(d3-2.0)) > 1e-6 {
		t.Errorf("expected 2.0 for opposite vectors, got %f", d3)
	}
}

func TestSearchVectorFallback(t *testing.T) {
	m := NewManager(Config{BaseDataDir: t.TempDir(), MaxReaders: 1})
	defer m.CloseAll()

	ctxA := core.WithTenant(context.Background(), core.TenantContext{UID: "userA", WorkspaceID: "ws1"})
	ctxB := core.WithTenant(context.Background(), core.TenantContext{UID: "userB", WorkspaceID: "ws1"}) // another tenant

	dbA, err := m.GetDB(ctxA)
	if err != nil {
		t.Fatal(err)
	}

	dbB, err := m.GetDB(ctxB)
	if err != nil {
		t.Fatal(err)
	}

	// Insert into User A
	vec1, _ := Float32ToBytes([]float32{1, 0, 0})
	vec2, _ := Float32ToBytes([]float32{0.9, 0.1, 0}) // similar
	vec3, _ := Float32ToBytes([]float32{0, 1, 0})     // completely different

	err = dbA.Write(ctxA, func(tx *sql.Tx) error {
		q := "INSERT INTO memories (id, uid, workspace_id, content, embedding) VALUES (?, ?, ?, ?, ?)"
		tx.Exec(q, "m1", "userA", "ws1", "Memory 1", vec1)
		tx.Exec(q, "m2", "userA", "ws1", "Memory 2", vec2)
		tx.Exec(q, "m3", "userA", "ws1", "Memory 3", vec3)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Insert into User B to test isolation
	vec4, _ := Float32ToBytes([]float32{1, 0, 0.1})
	dbB.Write(ctxB, func(tx *sql.Tx) error {
		tx.Exec("INSERT INTO memories (id, uid, workspace_id, content, embedding) VALUES (?, ?, ?, ?, ?)", "m4", "userB", "ws1", "Memory 4 (UserB)", vec4)
		return nil
	})

	// Search as User A
	queryVec := []float32{1, 0, 0}
	results, err := dbA.SearchVectorFallback(ctxA, queryVec, 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results due to limit, got %d", len(results))
	}

	if results[0].ID != "m1" {
		t.Errorf("expected top result to be m1, got %s", results[0].ID)
	}
	if results[1].ID != "m2" {
		t.Errorf("expected second result to be m2, got %s", results[1].ID)
	}

	// Search as User B (Should only see m4 despite query being close to m1)
	resultsB, err := dbB.SearchVectorFallback(ctxB, queryVec, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(resultsB) != 1 || resultsB[0].ID != "m4" {
		t.Errorf("cross tenant vector leakage detected or incorrect result")
	}
}

func TestSearchFTS_SanitizesQuery(t *testing.T) {
	m := NewManager(Config{BaseDataDir: t.TempDir(), MaxReaders: 1})
	defer m.CloseAll()

	ctx := core.WithTenant(context.Background(), core.TenantContext{UID: "userA", WorkspaceID: "ws1"})
	dbA, err := m.GetDB(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = dbA.Write(ctx, func(tx *sql.Tx) error {
		q := "INSERT INTO memories (id, uid, workspace_id, content, metadata) VALUES (?, ?, ?, ?, ?)"
		_, err := tx.Exec(q, "m1", "userA", "ws1", "O usuário gosta de café preto e tecnologia Go.", "{}")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := dbA.SearchFTS(ctx, "café preto?", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatalf("expected at least one FTS result")
	}
}
