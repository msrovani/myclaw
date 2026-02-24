package observability

import (
	"sync"
	"testing"
	"time"
)

func TestMetrics_Counter(t *testing.T) {
	m := NewMetrics()
	c := m.Counter("requests_total")
	c.Add(1)
	c.Add(5)

	if got := c.Load(); got != 6 {
		t.Errorf("Counter = %d, want 6", got)
	}

	// Same name returns same counter
	c2 := m.Counter("requests_total")
	if c2.Load() != 6 {
		t.Error("Counter should return same instance for same name")
	}
}

func TestMetrics_Gauge(t *testing.T) {
	m := NewMetrics()
	g := m.Gauge("goroutines")
	g.Store(42)

	if got := g.Load(); got != 42 {
		t.Errorf("Gauge = %d, want 42", got)
	}

	g.Store(10)
	if got := g.Load(); got != 10 {
		t.Errorf("Gauge = %d, want 10", got)
	}
}

func TestMetrics_Histogram(t *testing.T) {
	m := NewMetrics()
	m.RecordLatency("db_query", 100*time.Millisecond)
	m.RecordLatency("db_query", 200*time.Millisecond)
	m.RecordLatency("db_query", 150*time.Millisecond)

	snap := m.Snapshot()
	h, ok := snap.Histograms["db_query"]
	if !ok {
		t.Fatal("histogram db_query not found")
	}

	if h.Count != 3 {
		t.Errorf("Count = %d, want 3", h.Count)
	}
	if h.Min != int64(100*time.Millisecond) {
		t.Errorf("Min = %d, want %d", h.Min, int64(100*time.Millisecond))
	}
	if h.Max != int64(200*time.Millisecond) {
		t.Errorf("Max = %d, want %d", h.Max, int64(200*time.Millisecond))
	}
}

func TestMetrics_Snapshot(t *testing.T) {
	m := NewMetrics()
	m.Counter("a").Add(10)
	m.Gauge("b").Store(20)
	m.RecordLatency("c", time.Second)

	snap := m.Snapshot()
	if snap.Counters["a"] != 10 {
		t.Errorf("counter a = %d, want 10", snap.Counters["a"])
	}
	if snap.Gauges["b"] != 20 {
		t.Errorf("gauge b = %d, want 20", snap.Gauges["b"])
	}
	if snap.Histograms["c"].Count != 1 {
		t.Errorf("histogram c count = %d, want 1", snap.Histograms["c"].Count)
	}
}

func TestMetrics_ConcurrentAccess(t *testing.T) {
	m := NewMetrics()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Counter("concurrent").Add(1)
			m.Gauge("concurrent_g").Store(1)
			m.RecordLatency("concurrent_h", time.Millisecond)
			_ = m.Snapshot()
		}()
	}

	wg.Wait()

	if got := m.Counter("concurrent").Load(); got != 100 {
		t.Errorf("concurrent counter = %d, want 100", got)
	}
}
