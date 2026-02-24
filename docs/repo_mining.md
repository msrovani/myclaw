# XXXCLAW — Repo Mining Report

## Date

2026-02-23

## Repositories Analyzed

| Repository | Status | Language | License | Key Patterns |
|---|---|---|---|---|
| `mem0ai/mem0` | ✅ Public | Python | Apache-2.0 | Multi-level memory, user/session/agent scopes, embedding retrieval |
| `mem0ai/mem0-mcp` | ✅ Public | Python | Apache-2.0 | MCP tool contract, JSON API for memory operations |
| `asg017/sqlite-vec` | ✅ Public | C | MIT + Apache-2.0 | vec0 virtual tables, Go bindings, KNN queries |
| `openclaw/openclaw` | ❌ 404 | — | — | — |
| `openclaw/clawhub` | ❌ 404 | — | — | — |
| `openclaw/skills` | ❌ 404 | — | — | — |
| `zeroclaw-labs/zeroclaw` | ❌ 404 | — | — | — |
| Workspace `.agent` | ✅ Local | Markdown | — | 20 agents, 539 skills, 11 workflows |

## Pattern Matrix — Portable to Go

### Memory API (from Mem0)

| Operation | Mem0 Python | XXXCLAW Go Target |
|---|---|---|
| `memory.add(messages, user_id)` | Store memories from conversations | `AddMemory(ctx, messages, opts)` |
| `memory.search(query, user_id, limit)` | Semantic search with filters | `SearchMemories(ctx, query, opts)` |
| `memory.get_all(user_id)` | List all memories for scope | `GetMemories(ctx, opts)` |
| `memory.update(memory_id, data)` | Update specific memory | `UpdateMemory(ctx, id, data)` |
| `memory.delete(memory_id)` | Delete specific memory | `DeleteMemory(ctx, id)` |
| `memory.delete_all(user_id)` | Clear all memories for scope | `DeleteAllMemories(ctx, opts)` |

### Memory Scopes (from Mem0)

| Scope | Mem0 Support | XXXCLAW Target |
|---|---|---|
| `user_id` | ✅ Primary scope | ✅ |
| `session_id` | ✅ Session isolation | ✅ |
| `agent_id` | ✅ Agent-specific | ✅ |
| `tenant_id` | ❌ Not in OSS | ✅ Multi-tenant from start |

### Vector Search (from sqlite-vec)

| Feature | sqlite-vec | XXXCLAW Usage |
|---|---|---|
| `vec0` virtual table | ✅ | Memory embedding storage |
| Float vectors | ✅ | Default for semantic search |
| Int8 vectors | ✅ | Compact mode (optional) |
| KNN queries (`match` + `order by distance`) | ✅ | Hybrid retrieval |
| Go bindings | ✅ `github.com/asg017/sqlite-vec/bindings/go` | Direct integration |
| Metadata columns | ✅ Auxiliary/partition keys | Scope filtering |

### MCP Bridge (from mem0-mcp)

| Tool | mem0-mcp | XXXCLAW Target |
|---|---|---|
| `add_memory` | ✅ | ✅ (optional MCP server mode) |
| `search_memories` | ✅ | ✅ |
| `get_memories` | ✅ | ✅ |
| `update_memory` | ✅ | ✅ |
| `delete_memory` | ✅ | ✅ |
| `list_entities` | ✅ | ✅ |

### Concurrency Patterns (from workspace skills)

| Pattern | Source | XXXCLAW Application |
|---|---|---|
| Worker pools | `go-concurrency-patterns` | Per-domain worker pools |
| Pipeline stages | `go-concurrency-patterns` | Embedding → retrieval → ranking |
| Fan-out/fan-in | `golang-pro` | Parallel provider queries |
| Context cancellation | `golang-pro` | Timeout all DB/LLM calls |
| errgroup | `golang-pro` | Structured concurrency |
| sync.Pool | `golang-pro` | Buffer reuse on hot paths |

## Key Architectural Decisions from Mining

1. **Go-native Mem0**: Port API contract, not implementation. Go idioms (interfaces, context, error returns).
2. **sqlite-vec over external vector DB**: Local-first, single binary, no network dependency.
3. **Serialized writer**: Inspired by SQLite concurrency best practices — single writer goroutine avoids lock storms.
4. **Layered memory**: Extend Mem0's 3 scopes with configurable layers (working → short → long → episodic → procedural → system).
5. **Provider abstraction**: Clean interface so Ollama/Gemini/Claude/DeepSeek are interchangeable adapters.
