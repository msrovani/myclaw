package skills

import (
	"context"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
)

// MockSkill implements a simple skill for testing.
type MockSkill struct {
	id     string
	panics bool
}

func (m MockSkill) ID() string          { return m.id }
func (m MockSkill) Description() string { return "Mock skill for testing" }

func (m MockSkill) Execute(ctx context.Context, req Request) (Response, error) {
	if m.panics {
		panic("boom")
	}

	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return Response{Error: "no context"}, err
	}

	return Response{
		Result: map[string]any{
			"echo_uid": tc.UID,
			"msg":      "success",
		},
	}, nil
}

func TestRegistry_Execution(t *testing.T) {
	r := NewRegistry()
	_ = r.Register(MockSkill{id: "mock_safe", panics: false})
	_ = r.Register(MockSkill{id: "mock_panic", panics: true})

	// 1. Test missing context (Deny by default rule)
	_, err := r.Execute(context.Background(), Request{SkillID: "mock_safe"})
	if err == nil {
		t.Fatal("Execute should fail without TenantContext")
	}

	// 2. Test successful isolated execution
	ctxA := core.WithTenant(context.Background(), core.TenantContext{UID: "userA", WorkspaceID: "ws1"})
	req := Request{SkillID: "mock_safe"}
	resp, err := r.Execute(ctxA, req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if resp.Result["echo_uid"] != "userA" {
		t.Errorf("Expected echo_uid 'userA', got %v", resp.Result["echo_uid"])
	}

	// 3. Test Panic Recovery
	respPanic, errPanic := r.Execute(ctxA, Request{SkillID: "mock_panic"})
	if errPanic == nil {
		t.Fatal("Expected panic error")
	}
	if respPanic.Error == "" {
		t.Fatal("Expected panic error in response struct")
	}
}

func TestRegistry_Duplicate(t *testing.T) {
	r := NewRegistry()
	r.Register(MockSkill{id: "mock"})
	err := r.Register(MockSkill{id: "mock"})
	if err == nil {
		t.Fatal("Expected error on duplicate skill registration")
	}
}
