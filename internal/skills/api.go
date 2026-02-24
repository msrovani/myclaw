package skills

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/msrovani/myclaw/internal/core"
)

// Request defines the input for a skill execution.
// Payload is typically a JSON-unmarshalled map.
type Request struct {
	SkillID string         `json:"skill_id"`
	Payload map[string]any `json:"payload"`
}

// Response defines the output of a skill execution.
type Response struct {
	Result map[string]any `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

// Skill represents an executable tool or workflow within the system.
type Skill interface {
	// ID returns the unique identifier for the skill (e.g., "search_web", "read_file").
	ID() string

	// Description provides the LLM with context on what this skill does.
	Description() string

	// Execute runs the skill. It MUST receive a valid TenantContext via the context.
	// This ensures that tools (like reading memory or files) operate within the tenant's exact permissions.
	Execute(ctx context.Context, req Request) (Response, error)
}

// ensureTenant enforces that every skill invocation runs under a strict isolation context.
func ensureTenant(ctx context.Context) (core.TenantContext, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return tc, fmt.Errorf("skill execution blocked: %w", err)
	}
	return tc, nil
}

// ParsePayload is a helper to extract a specific struct from the generic map Payload.
func ParsePayload(req Request, dest any) error {
	b, err := json.Marshal(req.Payload)
	if err != nil {
		return fmt.Errorf("skill payload marshal err: %w", err)
	}
	return json.Unmarshal(b, dest)
}
