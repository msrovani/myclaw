package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPoolConfig configures a domain-specific worker pool.
type WorkerPoolConfig struct {
	Name      string
	Workers   int
	QueueSize int
	OnPanic   func(r any)
}

// Task is a unit of work submitted to a worker pool.
type Task struct {
	Name    string
	Execute func(ctx context.Context) error
}

// WorkerPool runs tasks concurrently with a bounded queue and fixed workers.
// Provides backpressure when the queue is full — Submit blocks until space is available.
type WorkerPool struct {
	name    string
	tasks   chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	onPanic func(r any)

	// Metrics
	submitted atomic.Int64
	completed atomic.Int64
	failed    atomic.Int64
	queueSize int
}

// NewWorkerPool creates and starts a worker pool.
func NewWorkerPool(cfg WorkerPoolConfig) *WorkerPool {
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		name:      cfg.Name,
		tasks:     make(chan Task, cfg.QueueSize),
		ctx:       ctx,
		cancel:    cancel,
		onPanic:   cfg.OnPanic,
		queueSize: cfg.QueueSize,
	}

	for i := 0; i < cfg.Workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	slog.Info("worker pool started",
		"name", cfg.Name,
		"workers", cfg.Workers,
		"queue_size", cfg.QueueSize,
	)
	return wp
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			wp.executeTask(id, task)
		}
	}
}

func (wp *WorkerPool) executeTask(workerID int, task Task) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			wp.failed.Add(1)
			slog.Error("worker pool: task panicked",
				"pool", wp.name,
				"worker", workerID,
				"task", task.Name,
				"panic", fmt.Sprint(r),
			)
			if wp.onPanic != nil {
				wp.onPanic(r)
			}
		}
	}()

	if err := task.Execute(wp.ctx); err != nil {
		wp.failed.Add(1)
		slog.Warn("worker pool: task failed",
			"pool", wp.name,
			"worker", workerID,
			"task", task.Name,
			"error", err,
			"duration", time.Since(start),
		)
	} else {
		wp.completed.Add(1)
	}
}

// Submit adds a task to the pool. Blocks if the queue is full (backpressure).
// Returns error if the pool is shut down or context is cancelled.
func (wp *WorkerPool) Submit(ctx context.Context, task Task) error {
	wp.submitted.Add(1)

	if err := wp.ctx.Err(); err != nil {
		return fmt.Errorf("worker pool %q is shut down", wp.name)
	}

	select {
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool %q is shut down", wp.name)
	case <-ctx.Done():
		return ctx.Err()
	case wp.tasks <- task:
		return nil
	}
}

// TrySubmit adds a task without blocking. Returns false if the queue is full.
func (wp *WorkerPool) TrySubmit(task Task) bool {
	if err := wp.ctx.Err(); err != nil {
		return false
	}

	select {
	case <-wp.ctx.Done():
		return false
	case wp.tasks <- task:
		wp.submitted.Add(1)
		return true
	default:
		return false
	}
}

// Stats returns pool statistics.
type PoolStats struct {
	Name      string `json:"name"`
	QueueLen  int    `json:"queue_len"`
	QueueCap  int    `json:"queue_cap"`
	Submitted int64  `json:"submitted"`
	Completed int64  `json:"completed"`
	Failed    int64  `json:"failed"`
}

// Stats returns a snapshot of pool metrics.
func (wp *WorkerPool) Stats() PoolStats {
	return PoolStats{
		Name:      wp.name,
		QueueLen:  len(wp.tasks),
		QueueCap:  cap(wp.tasks),
		Submitted: wp.submitted.Load(),
		Completed: wp.completed.Load(),
		Failed:    wp.failed.Load(),
	}
}

// Close gracefully shuts down the worker pool.
// Stops accepting new tasks, drains the queue, waits for workers.
func (wp *WorkerPool) Close() error {
	wp.cancel()
	close(wp.tasks)
	wp.wg.Wait()
	slog.Info("worker pool stopped",
		"name", wp.name,
		"completed", wp.completed.Load(),
		"failed", wp.failed.Load(),
	)
	return nil
}
