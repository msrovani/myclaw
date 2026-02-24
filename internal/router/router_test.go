package router

import (
	"context"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
	"github.com/msrovani/myclaw/internal/providers"
)

// MockProvider satisfies the Provider interface for routing tests.
type MockProvider struct {
	id string
}

func (m MockProvider) ID() string { return m.id }

func (m MockProvider) Generate(ctx context.Context, req providers.GenerateRequest) (providers.GenerateResponse, error) {
	return providers.GenerateResponse{
		Content:      "Mock string for " + m.id,
		InputTokens:  10,
		OutputTokens: 20,
		CostUSD:      0.015,
	}, nil
}

func (m MockProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	return []float32{1, 0, 0}, nil
}

func TestRouter_And_Economy(t *testing.T) {
	mgr := db.NewManager(db.Config{BaseDataDir: t.TempDir()})
	defer mgr.CloseAll()

	economy := NewEconomy(mgr)

	// Register two mocks: a local and a complex cloud one.
	r := NewRouter([]providers.Provider{
		MockProvider{id: "ollama"},
		MockProvider{id: "claude"},
	}, economy)

	ctxA := core.WithTenant(context.Background(), core.TenantContext{UID: "u1", WorkspaceID: "w1"})
	ctxB := core.WithTenant(context.Background(), core.TenantContext{UID: "u2", WorkspaceID: "w1"})

	// Test 1: Route to local due to LowLatencyOnly
	req := providers.GenerateRequest{Model: "qwen2.5"}
	resp, err := r.Route(ctxA, req, Policy{LowLatencyOnly: true})

	if err != nil {
		t.Fatalf("route failed: %v", err)
	}
	if resp.Content != "Mock string for ollama" {
		t.Errorf("expected ollama router resolution, got %s", resp.Content)
	}

	// The routing should have transparently charged User A 0.015 USD in tokens.
	costA, _ := economy.GetTotalCost(ctxA)
	if costA != 0.015 {
		t.Errorf("Economy tracking failed for User A! Expected 0.015, got %v", costA)
	}

	// Test 2: User B should have exactly 0.0 spent (Isolaton proof of Token Economy)
	costB, _ := economy.GetTotalCost(ctxB)
	if costB != 0.0 {
		t.Errorf("Cross-tenant economy leakage! User B has %v, expected 0", costB)
	}

	// Test 3: Complex routing
	resp, _ = r.Route(ctxB, req, Policy{ComplexReasoning: true})
	if resp.Content != "Mock string for claude" {
		t.Errorf("expected claude resolution, got %s", resp.Content)
	}

	// User B should now be billed.
	costBAfter, _ := economy.GetTotalCost(ctxB)
	if costBAfter != 0.015 {
		t.Errorf("Economy tracking failed for User B! Expected 0.015, got %v", costBAfter)
	}
}
