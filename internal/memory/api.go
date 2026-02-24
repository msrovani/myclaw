package memory

import (
	"context"
	"time"
)

// Layer represents the cognitive depth of a memory.
type Layer string

const (
	LayerWorking    Layer = "working"    // Current context window, immediate
	LayerShortTerm  Layer = "short_term" // Session-bound
	LayerLongTerm   Layer = "long_term"  // General semantic knowledge
	LayerEpisodic   Layer = "episodic"   // Specific past events
	LayerProcedural Layer = "procedural" // Implicit skills, rules, preferences
	LayerSystem     Layer = "system"     // Core agent directives (system prompts)
)

// Memory represents a single unit of knowledge.
type Memory struct {
	ID        string         `json:"id"`
	Content   string         `json:"content"`
	Layer     Layer          `json:"layer"`
	SessionID string         `json:"session_id,omitempty"`
	AgentID   string         `json:"agent_id,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Embedding []float32      `json:"embedding,omitempty"`
	Score     float32        `json:"score,omitempty"` // populated during search
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// AddOptions configures the addition of a memory.
type AddOptions struct {
	SessionID string
	AgentID   string
	Metadata  map[string]any
	Layer     Layer // Defaults to LongTerm if unspecified
}

// SearchOptions configures hybrid retrieval.
type SearchOptions struct {
	SessionID string
	AgentID   string
	Layers    []Layer
	Limit     int
	// Hybrid search controls
	Threshold float32
}

// API defines the Mem0-compatible persistent memory contract for XXXCLAW.
// ALL methods MUST receive a context containing a valid core.TenantContext.
type API interface {
	// AddMemory stores new knowledge into the tenant's isolated storage.
	AddMemory(ctx context.Context, content string, opts AddOptions) (Memory, error)

	// SearchMemories performs hybrid retrieval (Vector + FTS) for the tenant.
	SearchMemories(ctx context.Context, query string, queryEmbedding []float32, opts SearchOptions) ([]Memory, error)

	// GetMemory retrieves a specific memory by ID (strictly isolated to the tenant).
	GetMemory(ctx context.Context, id string) (Memory, error)

	// UpdateMemory modifies an existing memory.
	UpdateMemory(ctx context.Context, id string, content string, metadata map[string]any) (Memory, error)

	// DeleteMemory removes a memory.
	DeleteMemory(ctx context.Context, id string) error
}
