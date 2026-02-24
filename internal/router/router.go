package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/msrovani/myclaw/internal/providers"
)

// Policy defines the parameters for model routing.
type Policy struct {
	MaxCostAllowed   float64
	LowLatencyOnly   bool // forces local model
	ComplexReasoning bool // forces advanced model
}

// Router orchestrates provider selection out of a pool of adapters.
type Router struct {
	adapters map[string]providers.Provider
	economy  *Economy
}

// NewRouter creates a new LLM provider router.
func NewRouter(adapters []providers.Provider, economy *Economy) *Router {
	r := &Router{
		adapters: make(map[string]providers.Provider),
		economy:  economy,
	}
	for _, a := range adapters {
		r.adapters[a.ID()] = a
	}
	return r
}

// Route executes an LLM generation based on policy and available providers.
func (r *Router) Route(ctx context.Context, req providers.GenerateRequest, pol Policy) (providers.GenerateResponse, error) {
	// Simple routing heuristic
	var selected string

	if pol.LowLatencyOnly {
		selected = "ollama" // Always prefer local if latency/privacy is key
	} else if pol.ComplexReasoning {
		// Prefer Claude or DeepSeek for complexity
		if _, ok := r.adapters["claude"]; ok {
			selected = "claude"
		} else if _, ok := r.adapters["deepseek"]; ok {
			selected = "deepseek"
		} else {
			selected = "gemini" // fallback
		}
	} else {
		// Default to Gemini or whatever is registered
		for id := range r.adapters {
			selected = id
			break
		}
	}

	adapter, ok := r.adapters[selected]
	if !ok {
		return providers.GenerateResponse{}, fmt.Errorf("router: no active provider matching policy (wanted %s)", selected)
	}

	// Intercept and Execute
	resp, err := adapter.Generate(ctx, req)
	if err != nil {
		slog.Error("router: provider generation failed", "provider", selected, "error", err)
		return resp, err
	}

	// Record token usage if economy is active
	if r.economy != nil {
		err = r.economy.RecordUsage(ctx, adapter.ID(), req.Model, resp.InputTokens, resp.OutputTokens, resp.CostUSD)
		if err != nil {
			slog.Warn("router: failed to record token usage", "error", err)
		}
	}

	return resp, nil
}
