package core

import (
	"context"
	"log/slog"
	"sync"
)

// Event represents a typed event in the system.
type Event struct {
	Topic   string
	Payload any
}

// Handler processes events. Must be safe for concurrent use.
type Handler func(ctx context.Context, e Event)

// EventBus provides a typed, in-process pub/sub event bus.
// Subscriptions are topic-based. Handlers run concurrently via goroutines.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	wg       sync.WaitGroup
	closed   chan struct{}
}

// NewEventBus creates a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		closed:   make(chan struct{}),
	}
}

// Subscribe registers a handler for a topic. Thread-safe.
func (eb *EventBus) Subscribe(topic string, h Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[topic] = append(eb.handlers[topic], h)
}

// Publish sends an event to all handlers registered for the topic.
// Each handler runs in its own goroutine. Non-blocking for the publisher.
// Returns false if the bus is closed.
func (eb *EventBus) Publish(ctx context.Context, e Event) bool {
	select {
	case <-eb.closed:
		return false
	default:
	}

	eb.mu.RLock()
	handlers := eb.handlers[e.Topic]
	eb.mu.RUnlock()

	for _, h := range handlers {
		eb.wg.Add(1)
		go func(handler Handler) {
			defer eb.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					slog.Error("eventbus: handler panicked",
						"topic", e.Topic,
						"panic", r,
					)
				}
			}()
			handler(ctx, e)
		}(h)
	}
	return true
}

// Close shuts down the event bus and waits for all in-flight handlers.
func (eb *EventBus) Close() error {
	select {
	case <-eb.closed:
		return nil // already closed
	default:
		close(eb.closed)
	}
	eb.wg.Wait()
	return nil
}
