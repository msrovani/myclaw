package observability

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics provides a thread-safe metrics registry for tracking
// counters, gauges, and histograms across all XXXCLAW subsystems.
type Metrics struct {
	mu       sync.RWMutex
	counters map[string]*atomic.Int64
	gauges   map[string]*atomic.Int64
	histMu   sync.Mutex
	hists    map[string]*Histogram
}

// Histogram tracks latency distributions with min/max/avg/count.
type Histogram struct {
	mu    sync.Mutex
	count int64
	sum   int64 // nanoseconds
	min   int64
	max   int64
}

// NewMetrics creates a new metrics registry.
func NewMetrics() *Metrics {
	return &Metrics{
		counters: make(map[string]*atomic.Int64),
		gauges:   make(map[string]*atomic.Int64),
		hists:    make(map[string]*Histogram),
	}
}

// Counter returns or creates a named counter.
func (m *Metrics) Counter(name string) *atomic.Int64 {
	m.mu.RLock()
	c, ok := m.counters[name]
	m.mu.RUnlock()
	if ok {
		return c
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// Double-check after acquiring write lock.
	if c, ok = m.counters[name]; ok {
		return c
	}
	c = &atomic.Int64{}
	m.counters[name] = c
	return c
}

// Gauge returns or creates a named gauge.
func (m *Metrics) Gauge(name string) *atomic.Int64 {
	m.mu.RLock()
	g, ok := m.gauges[name]
	m.mu.RUnlock()
	if ok {
		return g
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if g, ok = m.gauges[name]; ok {
		return g
	}
	g = &atomic.Int64{}
	m.gauges[name] = g
	return g
}

// RecordLatency records a duration in the named histogram.
func (m *Metrics) RecordLatency(name string, d time.Duration) {
	m.histMu.Lock()
	h, ok := m.hists[name]
	if !ok {
		h = &Histogram{min: int64(d)}
		m.hists[name] = h
	}
	m.histMu.Unlock()

	ns := int64(d)

	h.mu.Lock()
	defer h.mu.Unlock()
	h.count++
	h.sum += ns
	if ns < h.min || h.count == 1 {
		h.min = ns
	}
	if ns > h.max {
		h.max = ns
	}
}

// Snapshot returns a point-in-time copy of all metrics.
type Snapshot struct {
	Counters   map[string]int64             `json:"counters"`
	Gauges     map[string]int64             `json:"gauges"`
	Histograms map[string]HistogramSnapshot `json:"histograms"`
}

// HistogramSnapshot is a point-in-time view of a histogram.
type HistogramSnapshot struct {
	Count int64   `json:"count"`
	Sum   int64   `json:"sum_ns"`
	Min   int64   `json:"min_ns"`
	Max   int64   `json:"max_ns"`
	Avg   float64 `json:"avg_ns"`
}

// Snapshot returns a point-in-time copy of all metrics.
func (m *Metrics) Snapshot() Snapshot {
	s := Snapshot{
		Counters:   make(map[string]int64),
		Gauges:     make(map[string]int64),
		Histograms: make(map[string]HistogramSnapshot),
	}

	m.mu.RLock()
	for k, v := range m.counters {
		s.Counters[k] = v.Load()
	}
	for k, v := range m.gauges {
		s.Gauges[k] = v.Load()
	}
	m.mu.RUnlock()

	m.histMu.Lock()
	for k, h := range m.hists {
		h.mu.Lock()
		hs := HistogramSnapshot{
			Count: h.count,
			Sum:   h.sum,
			Min:   h.min,
			Max:   h.max,
		}
		if h.count > 0 {
			hs.Avg = float64(h.sum) / float64(h.count)
		}
		h.mu.Unlock()
		s.Histograms[k] = hs
	}
	m.histMu.Unlock()

	return s
}
