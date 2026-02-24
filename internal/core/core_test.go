package core

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// --- EventBus Tests ---

func TestEventBus_PubSub(t *testing.T) {
	eb := NewEventBus()
	defer eb.Close()

	var received atomic.Int64
	eb.Subscribe("test.topic", func(ctx context.Context, e Event) {
		received.Add(1)
	})

	eb.Publish(context.Background(), Event{Topic: "test.topic", Payload: "hello"})
	time.Sleep(50 * time.Millisecond)

	if got := received.Load(); got != 1 {
		t.Errorf("received = %d, want 1", got)
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	eb := NewEventBus()
	defer eb.Close()

	var count atomic.Int64
	for i := 0; i < 5; i++ {
		eb.Subscribe("multi", func(ctx context.Context, e Event) {
			count.Add(1)
		})
	}

	eb.Publish(context.Background(), Event{Topic: "multi"})
	time.Sleep(50 * time.Millisecond)

	if got := count.Load(); got != 5 {
		t.Errorf("handlers called = %d, want 5", got)
	}
}

func TestEventBus_NoSubscribers(t *testing.T) {
	eb := NewEventBus()
	defer eb.Close()

	// Should not panic with no subscribers.
	ok := eb.Publish(context.Background(), Event{Topic: "empty"})
	if !ok {
		t.Error("Publish should return true for open bus")
	}
}

func TestEventBus_HandlerPanicRecovery(t *testing.T) {
	eb := NewEventBus()
	defer eb.Close()

	var afterPanic atomic.Int64
	eb.Subscribe("panic", func(ctx context.Context, e Event) {
		panic("test panic")
	})
	eb.Subscribe("after", func(ctx context.Context, e Event) {
		afterPanic.Add(1)
	})

	eb.Publish(context.Background(), Event{Topic: "panic"})
	eb.Publish(context.Background(), Event{Topic: "after"})
	time.Sleep(50 * time.Millisecond)

	if got := afterPanic.Load(); got != 1 {
		t.Errorf("after panic handler = %d, want 1", got)
	}
}

func TestEventBus_Close_RejectsPublish(t *testing.T) {
	eb := NewEventBus()
	eb.Close()

	ok := eb.Publish(context.Background(), Event{Topic: "t"})
	if ok {
		t.Error("Publish should return false after Close")
	}
}

func TestEventBus_Close_Idempotent(t *testing.T) {
	eb := NewEventBus()
	eb.Close()
	eb.Close() // should not panic
}

// --- WorkerPool Tests ---

func TestWorkerPool_BasicExecution(t *testing.T) {
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "test",
		Workers:   2,
		QueueSize: 10,
	})

	var done atomic.Int64
	for i := 0; i < 5; i++ {
		wp.Submit(context.Background(), Task{
			Name: fmt.Sprintf("task-%d", i),
			Execute: func(ctx context.Context) error {
				done.Add(1)
				return nil
			},
		})
	}

	time.Sleep(100 * time.Millisecond)
	wp.Close()

	if got := done.Load(); got != 5 {
		t.Errorf("completed = %d, want 5", got)
	}

	stats := wp.Stats()
	if stats.Completed != 5 {
		t.Errorf("stats.Completed = %d, want 5", stats.Completed)
	}
}

func TestWorkerPool_ErrorCounting(t *testing.T) {
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "err-test",
		Workers:   1,
		QueueSize: 10,
	})

	wp.Submit(context.Background(), Task{
		Name: "fail",
		Execute: func(ctx context.Context) error {
			return errors.New("failed")
		},
	})
	wp.Submit(context.Background(), Task{
		Name: "ok",
		Execute: func(ctx context.Context) error {
			return nil
		},
	})

	time.Sleep(100 * time.Millisecond)
	wp.Close()

	stats := wp.Stats()
	if stats.Failed != 1 {
		t.Errorf("stats.Failed = %d, want 1", stats.Failed)
	}
	if stats.Completed != 1 {
		t.Errorf("stats.Completed = %d, want 1", stats.Completed)
	}
}

func TestWorkerPool_PanicRecovery(t *testing.T) {
	var panicCaught atomic.Int64
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "panic-test",
		Workers:   1,
		QueueSize: 10,
		OnPanic: func(r any) {
			panicCaught.Add(1)
		},
	})

	wp.Submit(context.Background(), Task{
		Name: "panic",
		Execute: func(ctx context.Context) error {
			panic("boom")
		},
	})

	time.Sleep(100 * time.Millisecond)
	wp.Close()

	if got := panicCaught.Load(); got != 1 {
		t.Errorf("panics caught = %d, want 1", got)
	}
}

func TestWorkerPool_Backpressure(t *testing.T) {
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "bp-test",
		Workers:   1,
		QueueSize: 2,
	})

	blocker := make(chan struct{})
	// Fill the worker and queue.
	wp.Submit(context.Background(), Task{
		Name: "block",
		Execute: func(ctx context.Context) error {
			<-blocker
			return nil
		},
	})
	wp.Submit(context.Background(), Task{Name: "q1", Execute: func(ctx context.Context) error { return nil }})
	wp.Submit(context.Background(), Task{Name: "q2", Execute: func(ctx context.Context) error { return nil }})

	// TrySubmit should fail — queue is full.
	ok := wp.TrySubmit(Task{Name: "overflow", Execute: func(ctx context.Context) error { return nil }})
	if ok {
		t.Error("TrySubmit should return false when queue is full")
	}

	close(blocker)
	wp.Close()
}

func TestWorkerPool_SubmitAfterClose(t *testing.T) {
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "closed-test",
		Workers:   1,
		QueueSize: 10,
	})
	wp.Close()

	err := wp.Submit(context.Background(), Task{
		Name:    "late",
		Execute: func(ctx context.Context) error { return nil },
	})
	if err == nil {
		t.Error("Submit after Close should return error")
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	wp := NewWorkerPool(WorkerPoolConfig{
		Name:      "ctx-test",
		Workers:   1,
		QueueSize: 1,
	})
	defer wp.Close()

	// Fill the worker and queue so Submit would block.
	blocker := make(chan struct{})
	wp.Submit(context.Background(), Task{
		Name: "block",
		Execute: func(ctx context.Context) error {
			<-blocker
			return nil
		},
	})
	wp.Submit(context.Background(), Task{Name: "q1", Execute: func(ctx context.Context) error { return nil }})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := wp.Submit(ctx, Task{Name: "timeout", Execute: func(ctx context.Context) error { return nil }})
	if err == nil {
		t.Error("Submit should return error when context is cancelled")
	}

	close(blocker)
}
