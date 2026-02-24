# XXXCLAW

Modular agent system in Go with multicore concurrency, Mem0-like persistent memory, SQLite vector search, operational dashboard, and token economy.

## Features (Planned)

- **Multicore Concurrency** — Worker pools, pipelines, backpressure, zero goroutine leaks
- **Persistent Memory** — Mem0-compatible API, layered memory (working → long-term), hybrid retrieval
- **SQLite + Vector Search** — WAL, serialized writer, sqlite-vec for semantic search
- **Token Economy** — Multi-provider routing (Ollama/Gemini/Claude/DeepSeek), semantic cache, prompt compression
- **Dashboard** — templ + HTMX + ECharts operational UI
- **Skills/Workflows** — Runtime with `.agent` integration

## Quick Start

```bash
# Build
make build

# Run (defaults to :8080)
make run

# Test
make test

# Test with race detector
make test-race

# Benchmarks
make bench
```

## Configuration

All configuration via environment variables:

| Variable | Default | Description |
|---|---|---|
| `XXXCLAW_HTTP_ADDR` | `:8080` | HTTP server address |
| `XXXCLAW_LOG_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `XXXCLAW_LOG_FORMAT` | `json` | Log format (json/text) |
| `XXXCLAW_PPROF_ENABLED` | `true` | Enable pprof endpoint |
| `XXXCLAW_PPROF_ADDR` | `:6060` | pprof server address |
| `XXXCLAW_DB_PATH` | `data/xxxclaw.db` | SQLite database path |
| `XXXCLAW_VECTOR_ENABLED` | `true` | Enable vector search |
| `XXXCLAW_VECTOR_DIM` | `384` | Vector dimension |
| `XXXCLAW_ENV` | `dev` | Environment (dev/prod/edge) |
| `XXXCLAW_OLLAMA_URL` | `http://localhost:11434` | Ollama API URL |

## Project Structure

```
cmd/xxxclaw/          Entry point
internal/
  config/             Config loader
  core/               Event bus, workers, scheduler
  db/                 SQLite, migrations, vec
  memory/             Mem0-like engine
  providers/          LLM providers
  router/             Token economy + routing
  skills/             Skill runtime
  workflows/          Workflow runtime
  dashboard/          Dashboard backend
  observability/      Logging, metrics, pprof
web/                  Dashboard frontend
migrations/           SQLite migrations
docs/                 Documentation
```

## License

See [docs/license_matrix.md](docs/license_matrix.md) for dependency licenses.
