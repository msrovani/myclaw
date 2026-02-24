# XXXCLAW — Architecture Decision Records

## ADR-001: Go Monorepo with `internal/` Layout

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need a modular, portable project structure for a single deployable binary with clear module boundaries.

**Decision**: Single Go module with `cmd/xxxclaw/` entry point and `internal/` for all packages. Modules: `config`, `core`, `db`, `memory`, `providers`, `router`, `skills`, `workflows`, `dashboard`, `observability`.

**Rationale**: Go `internal/` prevents external imports, enforcing encapsulation. Single binary simplifies deployment. Monorepo avoids multi-module dependency management overhead.

---

## ADR-002: SQLite + WAL as Primary Storage

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need persistent local storage for memories, embeddings, metadata, and ops data. Must support concurrent reads with minimal lock contention.

**Decision**: SQLite3 with WAL mode enabled, serialized writer (single goroutine), multiple reader connections. Use `github.com/mattn/go-sqlite3` (CGo).

**Rationale**: Single-file database, zero server dependency, WAL enables concurrent reads during writes. Serialized writer avoids SQLITE_BUSY lock storms.

---

## ADR-003: sqlite-vec for Vector Search

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need local vector search for semantic memory retrieval without external services.

**Decision**: Use `sqlite-vec` (`github.com/asg017/sqlite-vec/bindings/go`) with `vec0` virtual tables. Feature flag to disable vector search (graceful degradation).

**Rationale**: Runs anywhere SQLite runs. Pure C, no dependencies. Go bindings available. Pre-v1 — isolate with feature flag for breaking changes.

---

## ADR-004: Mem0-Compatible Memory API

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need a persistent memory engine with multi-scope, multi-layer architecture.

**Decision**: Go-native implementation inspired by Mem0's API contract (`Add/Search/Get/Update/Delete/ListEntities`). Scopes: user, session, agent, tenant. Minimum 6 memory layers (working → system).

**Rationale**: Mem0's API is well-designed and widely adopted. Go-native implementation avoids Python dependency. Multi-tenant from start. Layers are configurable, not hardcoded.

---

## ADR-005: Provider Abstraction with Routing Policy

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need to support multiple LLM providers (local + cloud) with cost optimization.

**Decision**: `Provider` interface with adapters (Ollama, Gemini, Claude, DeepSeek). Routing policy engine decides per request based on cost, latency, privacy, complexity, budget.

**Rationale**: Decouples business logic from specific providers. Enables local-first strategy with cloud fallback. Token accounting per request.

---

## ADR-006: templ + HTMX + ECharts Dashboard

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need an operational dashboard without heavy JS frameworks.

**Decision**: Server-rendered HTML with `templ` (Go templates), `HTMX` for partial updates, `Alpine.js` for minimal client state, `ECharts` for charts.

**Rationale**: Minimal JS dependencies. Server-rendered = fast initial load. HTMX enables real-time updates without SPA complexity. ECharts handles complex visualizations.

---

## ADR-007: Structured Logging with slog

**Status**: Accepted  
**Date**: 2026-02-23

**Context**: Need structured, leveled logging with JSON output for production.

**Decision**: Use Go 1.21+ `log/slog` with JSON handler. Sensitive data redaction middleware.

**Rationale**: Standard library, zero dependencies, structured output, configurable levels.
