---
description: "Phase 5 — Build LLM provider router, token economy, semantic cache."
---

# 40 — Provider Router

## Steps

1. Define `Provider` interface (Chat, Embed, Complete).
2. Implement adapters: Ollama, Gemini, Claude, DeepSeek.
3. Build routing policy engine (cost, latency, privacy, complexity).
4. Implement token accounting (in/out, cost, per-request).
5. Build semantic cache (embedding-based response reuse).
6. Implement prompt compression pipeline.
7. Build fallback chain (local → cloud).
8. Measure token savings vs baseline.

## Agents Used

- `llm-routing-token-economy`
- `go-concurrency-perf`

## Outputs

- `internal/providers/` package
- `internal/router/` package
- Token accounting middleware
- Semantic cache
- Per-request metrics
