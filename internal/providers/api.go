package providers

import (
	"context"
)

// Role defines the message sender role.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message represents a single chat turn.
type Message struct {
	Role    Role
	Content string
	Name    string // optional, e.g. for tool responses
}

// GenerateRequest represents the standard input for any LLM provider.
type GenerateRequest struct {
	Messages    []Message
	Model       string
	Temperature float32
	MaxTokens   int
	System      string
}

// GenerateResponse represents the standardized output from any LLM provider.
type GenerateResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
	StopReason   string
	CostUSD      float64 // If known/applicable
}

// Provider defines the standard capability of an LLM.
// ALL methods MUST receive a valid core.TenantContext.
// This ensures that providers can optionally be tied to specific isolated API keys
// or accounting.
type Provider interface {
	// ID returns the provider slug (e.g., "ollama", "gemini").
	ID() string

	// Generate text from a sequence of messages.
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)

	// Embeddings generates float32 vectors for text.
	Embed(ctx context.Context, text string) ([]float32, error)
}
