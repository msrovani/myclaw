# XXXCLAW — Agent Index

## XXXCLAW-Specific Agents

| Agent | Purpose | When to Use | Inputs | Outputs |
|---|---|---|---|---|
| `architect-go-systems` | Go module layout, interfaces, DI, concurrency topology | New modules, interface design, ADRs | Requirements, constraints | Module layout, interfaces, ADR |
| `go-concurrency-perf` | Worker pools, pipelines, pprof, race detection, GC | Concurrency design, profiling, benchmarks | Bottleneck data, targets | Workers, pipelines, profiles |
| `memory-engineer-mem0like` | Layered memory, Mem0 API, hybrid retrieval, compaction | Memory schema, retrieval, policies | Scope/layer requirements | Memory engine, tests |
| `sqlite-vector-engineer` | SQLite WAL, writer queue, readers pool, sqlite-vec, FTS5 | Schema, migrations, vector integration | Data model, concurrency | DB layer, benchmarks |
| `llm-routing-token-economy` | Provider abstraction, routing, token accounting, cache | Provider adapters, budget, compression | API specs, policies | Providers, router, metrics |
| `dashboard-ops-ui` | templ+HTMX+ECharts operational dashboard | Dashboard panels, admin controls | Metrics sources, panel specs | Handlers, templates, charts |
| `qa-security-release` | Tests, benchmarks, audit, hardening, release | Testing, security review, release prep | Code, requirements | Tests, audits, artifacts |

## Pre-existing Agents (from workspace)

| Agent | Purpose |
|---|---|
| `orchestrator` | Multi-agent coordination |
| `project-planner` | Project planning |
| `backend-specialist` | Backend development |
| `database-architect` | Database design |
| `debugger` | Debugging |
| `performance-optimizer` | Performance optimization |
| `security-auditor` | Security auditing |
| `test-engineer` | Test engineering |

## XXXCLAW Workflows

| Workflow | Phase | Description |
|---|---|---|
| `/00_discovery_repo_mining` | 0 | Repo study, pattern extraction |
| `/10_architecture_design` | 1 | Go architecture, ADRs, scaffold |
| `/20_db_memory_schema` | 3 | SQLite schema, memory tables |
| `/30_backend_core` | 2 | Event bus, workers, scheduler |
| `/40_provider_router` | 5 | LLM providers, token economy |
| `/50_dashboard` | 7 | Operational dashboard |
| `/60_perf_tuning` | 8 | Profiling, optimization |
| `/70_hardening_release` | 9 | Security, release engineering |

## Key Skills (from 539 available)

| Skill | Relevance |
|---|---|
| `golang-pro` | Core Go development |
| `go-concurrency-patterns` | Concurrency patterns |
| `agent-memory-systems` | Memory architecture |
| `vector-database-engineer` | Vector search |
| `clean-code` | Code quality |
| `architecture-decision-records` | ADR writing |
| `observability-engineer` | Metrics, tracing |
