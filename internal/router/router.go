package router

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"sync"

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

	cacheMu sync.RWMutex
	cache   map[string]providers.GenerateResponse
}

// NewRouter creates a new LLM provider router.
func NewRouter(adapters []providers.Provider, economy *Economy) *Router {
	r := &Router{
		adapters: make(map[string]providers.Provider),
		economy:  economy,
		cache:    make(map[string]providers.GenerateResponse),
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

	// 0. Prompt Compression (Token savings)
	compressedMsgs := compressMessages(req.Messages)
	req.Messages = compressedMsgs

	// 1. Semantic/Exact Match Cache Check
	promptKey := generateCacheKey(req.Messages)
	r.cacheMu.RLock()
	cachedResp, found := r.cache[promptKey]
	r.cacheMu.RUnlock()

	if found {
		slog.Info("router: cache hit", "prompt_hash", promptKey[:8])
		// Cost is 0 on cache hit
		cachedResp.CostUSD = 0
		cachedResp.InputTokens = 0
		cachedResp.OutputTokens = 0
		return cachedResp, nil
	}

	// 2. Intercept and Execute
	resp, err := adapter.Generate(ctx, req)
	if err != nil {
		slog.Error("router: provider generation failed", "provider", selected, "error", err)
		return resp, err
	}

	// 3. Store in Cache
	r.cacheMu.Lock()
	r.cache[promptKey] = resp
	r.cacheMu.Unlock()

	// 4. Record token usage if economy is active
	if r.economy != nil {
		err = r.economy.RecordUsage(ctx, adapter.ID(), req.Model, resp.InputTokens, resp.OutputTokens, resp.CostUSD)
		if err != nil {
			slog.Warn("router: failed to record token usage", "error", err)
		}
	}

	return resp, nil
}

func generateCacheKey(msgs []providers.Message) string {
	var sb strings.Builder
	for _, m := range msgs {
		sb.WriteString(string(m.Role))
		sb.WriteString(":")
		sb.WriteString(m.Content)
		sb.WriteString("|")
	}
	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

// compressMessages provides a simplistic prompt compression by stripping
// redundant whitespaces, newlines, and truncating absurdly long messages.
// A full implementation would use a small local NLP model (e.g. LLMLingua).
func compressMessages(msgs []providers.Message) []providers.Message {
	compressed := make([]providers.Message, 0, len(msgs))
	for _, m := range msgs {
		content := m.Content
		// Strip redundant whitespaces
		content = strings.Join(strings.Fields(content), " ")

		// Truncate if insanely long (e.g., > 100k chars ~ 25k tokens), just to protect budget
		const maxChars = 100000
		if len(content) > maxChars {
			content = content[:maxChars] + "... [TRUNCATED BY ROUTER]"
		}

		compressed = append(compressed, providers.Message{
			Role:    m.Role,
			Content: content,
			Name:    m.Name,
		})
	}
	return compressed
}
