---
description: "Phase 8 — Profile, optimize, tune GC/pools/SQLite/retrieval/routing."
---

# 60 — Performance Tuning

## Steps

1. Run pprof profiles (CPU, heap, goroutine, mutex, block).
2. Identify top allocation sites and reduce with `sync.Pool`.
3. Tune GC parameters (`GOGC`, `GOMEMLIMIT`).
4. Tune SQLite PRAGMAs for production.
5. Optimize retrieval pipeline (batch, caching, index tuning).
6. Tune LLM routing (latency thresholds, cache sizes).
7. Run benchmarks before/after each optimization.
8. Document all changes with metrics in `docs/progress/`.

## Agents Used

- `go-concurrency-perf`
- `sqlite-vector-engineer`
- `llm-routing-token-economy`

## Outputs

- Profiling reports
- Optimization commits with benchmark diffs
- `docs/progress/perf_tuning.md`
