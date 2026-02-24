# XXXCLAW — Roadmap

## Phase 0 — Repo Mining & Plan ✅

- Study reference repos (Mem0, sqlite-vec, mem0-mcp)
- Generate docs: `repo_mining.md`, `license_matrix.md`, `architecture_decisions.md`
- Create `.agent` agents (7), workflows (8), `AGENT_INDEX.md`
- Project scaffold

## Phase 1 — Scaffold Base (Go)

- Go module, directory layout, Makefile
- Config loader (YAML + env)
- Structured logging (slog)
- Simple DI
- Health endpoint + pprof
- Observability base

## Phase 2 — Core Concurrency

- Event bus (typed, in-process)
- Bounded queues + backpressure
- Worker pools (I/O, embeddings, retrieval, tools, maintenance)
- Job scheduler
- `context.Context` + `errgroup`
- Queue/latency metrics
- `go test -race`

## Phase 3 — SQLite Core + sqlite-vec

- Schema + versioned migrations
- WAL + PRAGMAs
- Serialized writer queue
- Readers pool
- `vec0` integration
- FTS5 indices
- Lock/contention tests + benchmarks

## Phase 4 — Memory Engine (Mem0-like)

- Memory types/layers (6+ configurable)
- API: Add/Search/Get/Update/Delete/ListEntities
- Scopes: user/session/agent/tenant
- Hybrid retrieval + rank fusion
- Dedup, decay, compaction, retention

## Phase 5 — Provider Router + Token Economy

- Provider interface + adapters (Ollama, Gemini, Claude, DeepSeek)
- Routing policy engine
- Token accounting + budget
- Semantic cache
- Prompt compression
- Fallback chain + per-request metrics

## Phase 6 — Skills / Tools / Workflows Runtime

- Skill runtime + registry
- I/O contracts
- Concurrent execution with limits
- Audit + metrics
- `.agent` integration

## Phase 7 — Dashboard

- Admin API (Go net/http)
- Frontend (templ + HTMX + Alpine.js + ECharts)
- Panels: system, tokens, memory, router, skills, admin
- Responsive + dark mode

## Phase 8 — Hardening & Optimization

- pprof (CPU/heap/mutex/block)
- Allocation reduction, GC tuning
- SQLite tuning, retrieval tuning, routing tuning
- Before/after benchmarks

## Phase 9 — Release Engineering

- Documentation, build scripts, Dockerfile
- Deploy examples, changelog
- Extension guide (skill/tool/provider/memory policy)
