---
description: "Phase 3 — Design SQLite schema, memory layer tables, migration system."
---

# 20 — DB & Memory Schema

## Steps

1. Design SQLite schema for memories, embeddings, metadata, ops.
2. Create versioned migration files in `migrations/`.
3. Configure WAL + PRAGMAs.
4. Implement serialized writer queue.
5. Set up readers pool.
6. Integrate sqlite-vec (`vec0` virtual tables).
7. Add FTS5 indices.
8. Test lock/contention patterns.
9. Benchmark ingest and retrieval.

## Agents Used

- `sqlite-vector-engineer`
- `memory-engineer-mem0like`

## Outputs

- `migrations/*.sql`
- `internal/db/` package
- Benchmark results
