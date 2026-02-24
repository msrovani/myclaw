// Package container provides a simple, type-safe dependency injection
// container for XXXCLAW. No reflection, no magic — just explicit
// constructor injection with lifecycle management.
package container

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Container holds all application dependencies with lifecycle management.
type Container struct {
	mu        sync.RWMutex
	services  map[string]any
	closers   []func() error
	closeOnce sync.Once
}

// New creates a new dependency injection container.
func New() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

// Register adds a service to the container by name.
// If the service implements io.Closer, it will be closed on container shutdown.
func (c *Container) Register(name string, svc any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = svc

	// Track closeable services for shutdown.
	type closer interface {
		Close() error
	}
	if cl, ok := svc.(closer); ok {
		c.closers = append(c.closers, cl.Close)
	}
}

// Get retrieves a service by name. Returns nil if not found.
func (c *Container) Get(name string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.services[name]
}

// MustGet retrieves a service by name, panics if not found.
func (c *Container) MustGet(name string) any {
	svc := c.Get(name)
	if svc == nil {
		panic(fmt.Sprintf("container: service %q not registered", name))
	}
	return svc
}

// Resolve is a generic helper to retrieve and type-assert a service.
func Resolve[T any](c *Container, name string) (T, error) {
	svc := c.Get(name)
	if svc == nil {
		var zero T
		return zero, fmt.Errorf("container: service %q not registered", name)
	}
	typed, ok := svc.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("container: service %q is %T, want %T", name, svc, zero)
	}
	return typed, nil
}

// Shutdown closes all registered services in reverse order.
// Safe to call multiple times.
func (c *Container) Shutdown(ctx context.Context) error {
	var firstErr error
	c.closeOnce.Do(func() {
		slog.Info("container: shutting down services", "count", len(c.closers))
		// Close in reverse registration order (LIFO).
		for i := len(c.closers) - 1; i >= 0; i-- {
			select {
			case <-ctx.Done():
				firstErr = ctx.Err()
				return
			default:
			}
			if err := c.closers[i](); err != nil {
				slog.Error("container: service close error", "error", err)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	})
	return firstErr
}
