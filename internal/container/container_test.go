package container

import (
	"context"
	"errors"
	"testing"
)

func TestContainer_RegisterAndGet(t *testing.T) {
	c := New()
	c.Register("config", "value123")

	got := c.Get("config")
	if got != "value123" {
		t.Errorf("Get = %v, want %v", got, "value123")
	}
}

func TestContainer_GetNotFound(t *testing.T) {
	c := New()
	if got := c.Get("missing"); got != nil {
		t.Errorf("Get(missing) = %v, want nil", got)
	}
}

func TestContainer_MustGet_Panics(t *testing.T) {
	c := New()
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic for missing service")
		}
	}()
	c.MustGet("missing")
}

func TestContainer_Resolve_Generic(t *testing.T) {
	c := New()
	c.Register("port", 8080)

	val, err := Resolve[int](c, "port")
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if val != 8080 {
		t.Errorf("Resolve = %d, want 8080", val)
	}
}

func TestContainer_Resolve_TypeMismatch(t *testing.T) {
	c := New()
	c.Register("port", "8080")

	_, err := Resolve[int](c, "port")
	if err == nil {
		t.Fatal("Resolve should error on type mismatch")
	}
}

func TestContainer_Resolve_NotFound(t *testing.T) {
	c := New()
	_, err := Resolve[int](c, "missing")
	if err == nil {
		t.Fatal("Resolve should error on missing service")
	}
}

type mockCloser struct {
	closed bool
	err    error
}

func (m *mockCloser) Close() error {
	m.closed = true
	return m.err
}

func TestContainer_Shutdown_ClosesServices(t *testing.T) {
	c := New()
	svc1 := &mockCloser{}
	svc2 := &mockCloser{}

	c.Register("svc1", svc1)
	c.Register("svc2", svc2)

	err := c.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Shutdown error: %v", err)
	}

	if !svc1.closed {
		t.Error("svc1 should be closed")
	}
	if !svc2.closed {
		t.Error("svc2 should be closed")
	}
}

func TestContainer_Shutdown_ReverseOrder(t *testing.T) {
	c := New()
	var order []string

	c.Register("first", &orderCloser{name: "first", order: &order})
	c.Register("second", &orderCloser{name: "second", order: &order})
	c.Register("third", &orderCloser{name: "third", order: &order})

	c.Shutdown(context.Background())

	if len(order) != 3 {
		t.Fatalf("expected 3 closes, got %d", len(order))
	}
	if order[0] != "third" || order[1] != "second" || order[2] != "first" {
		t.Errorf("shutdown order = %v, want [third second first]", order)
	}
}

type orderCloser struct {
	name  string
	order *[]string
}

func (o *orderCloser) Close() error {
	*o.order = append(*o.order, o.name)
	return nil
}

func TestContainer_Shutdown_ReportsFirstError(t *testing.T) {
	c := New()
	expected := errors.New("close failed")
	c.Register("bad", &mockCloser{err: expected})

	err := c.Shutdown(context.Background())
	if err == nil {
		t.Fatal("Shutdown should report error")
	}
	if err.Error() != expected.Error() {
		t.Errorf("error = %v, want %v", err, expected)
	}
}

func TestContainer_Shutdown_Idempotent(t *testing.T) {
	c := New()
	svc := &mockCloser{}
	c.Register("svc", svc)

	c.Shutdown(context.Background())
	c.Shutdown(context.Background())
	// Close should only be called once thanks to sync.Once
}
