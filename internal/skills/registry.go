package skills

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Registry holds the available skills in the system.
// While skills themselves might be global (code), their *execution* is
// strictly scoped by the TenantContext passed to Execute().
type Registry struct {
	mu     sync.RWMutex
	skills map[string]Skill
}

var ErrSkillAlreadyRegistered = errors.New("skill already registered")

// NewRegistry creates a new Skill registry.
func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]Skill),
	}
}

// Register adds a skill to the registry.
func (r *Registry) Register(s Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := s.ID()
	if _, exists := r.skills[id]; exists {
		return fmt.Errorf("%w: %s", ErrSkillAlreadyRegistered, id)
	}
	r.skills[id] = s
	slog.Debug("skills: registered", "id", id)
	return nil
}

// Execute looks up a skill and runs it, enforcing Tenant bounds and measuring execution.
// It tracks audit logs (via slog) and panic recovery.
func (r *Registry) Execute(ctx context.Context, req Request) (resp Response, err error) {
	// 1. Mandatory Tenant Check (Deny by default)
	tc, err := ensureTenant(ctx)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	r.mu.RLock()
	skill, exists := r.skills[req.SkillID]
	r.mu.RUnlock()

	if !exists {
		err := fmt.Errorf("skill %q not found", req.SkillID)
		return Response{Error: err.Error()}, err
	}

	start := time.Now()

	// 2. Audit Logging (Start)
	slog.Info("skills: executing",
		"uid", tc.UID,
		"workspace_id", tc.WorkspaceID,
		"skill_id", req.SkillID,
	)

	// Panic Recovery during execution
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("skill %q panicked: %v", req.SkillID, rec)
			slog.Error("skills: panic recovered",
				"uid", tc.UID,
				"workspace_id", tc.WorkspaceID,
				"skill_id", req.SkillID,
				"panic", rec,
			)
			resp = Response{Error: err.Error()}
		}

		// 3. Audit Logging (End / Metrics)
		duration := time.Since(start)
		slog.Info("skills: completed",
			"uid", tc.UID,
			"workspace_id", tc.WorkspaceID,
			"skill_id", req.SkillID,
			"duration_ms", duration.Milliseconds(),
			"success", err == nil,
		)
	}()

	// 4. Execution
	resp, err = skill.Execute(ctx, req)
	if err != nil && resp.Error == "" {
		resp.Error = err.Error()
	}

	return resp, err
}

// List returns all registered skills (useful for LLM prompting).
func (r *Registry) List() []Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Skill, 0, len(r.skills))
	for _, s := range r.skills {
		list = append(list, s)
	}
	return list
}
