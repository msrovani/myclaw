package observability

import (
	"runtime"
	"time"
)

// SystemStats holds basic runtime metrics for the dashboard.
type SystemStats struct {
	Goroutines int
	AllocBytes uint64
	SysBytes   uint64
	Uptime     string
}

var startTime = time.Now()

// CollectSystemStats gathers high-level application metrics.
func CollectSystemStats() SystemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime).Round(time.Second)

	return SystemStats{
		Goroutines: runtime.NumGoroutine(),
		AllocBytes: m.Alloc,
		SysBytes:   m.Sys,
		Uptime:     uptime.String(),
	}
}
