---
description: "Phase 1 — Design Go architecture, create ADRs, scaffold project structure."
---

# 10 — Architecture Design

## Steps

1. Define module layout (`cmd/`, `internal/`, `web/`, `migrations/`, `docs/`).
2. Define interface contracts between modules.
3. Design dependency direction graph.
4. Write initial ADRs in `docs/architecture_decisions.md`.
5. Create `go.mod`, `Makefile`, project scaffold.
6. Implement config loader, structured logging, health endpoint, pprof.

## Agents Used

- `architect-go-systems`
- `go-concurrency-perf`

## Outputs

- Project scaffold (directories + go.mod)
- `docs/architecture_decisions.md`
- `docs/roadmap.md`
- Config loader, slog, health endpoint
