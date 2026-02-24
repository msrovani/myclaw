package core

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Job represents a scheduled recurring job.
type Job struct {
	Name     string
	Interval time.Duration
	Execute  func(ctx context.Context) error
}

// Scheduler runs recurring jobs at configured intervals.
// Each job runs in its own goroutine with independent timing.
type Scheduler struct {
	mu     sync.Mutex
	jobs   []jobEntry
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type jobEntry struct {
	job  Job
	stop context.CancelFunc
}

// NewScheduler creates a new job scheduler.
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Schedule adds a recurring job. The job starts immediately.
func (s *Scheduler) Schedule(job Job) {
	if job.Interval <= 0 {
		slog.Error("scheduler: invalid interval", "job", job.Name, "interval", job.Interval)
		return
	}

	jobCtx, jobCancel := context.WithCancel(s.ctx)

	s.mu.Lock()
	s.jobs = append(s.jobs, jobEntry{job: job, stop: jobCancel})
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run(jobCtx, job)

	slog.Info("scheduler: job registered",
		"job", job.Name,
		"interval", job.Interval,
	)
}

func (s *Scheduler) run(ctx context.Context, job Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("scheduler: job panicked",
							"job", job.Name,
							"panic", r,
						)
					}
				}()

				start := time.Now()
				if err := job.Execute(ctx); err != nil {
					slog.Warn("scheduler: job failed",
						"job", job.Name,
						"error", err,
						"duration", time.Since(start),
					)
				}
			}()
		}
	}
}

// Close stops all scheduled jobs and waits for them to finish.
func (s *Scheduler) Close() error {
	s.cancel()
	s.wg.Wait()
	slog.Info("scheduler: all jobs stopped")
	return nil
}
