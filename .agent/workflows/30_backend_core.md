---
description: "Phase 2 — Build core concurrency: event bus, worker pools, scheduler."
---

# 30 — Backend Core

## Steps

1. Implement typed event bus (publish/subscribe).
2. Build bounded queues with backpressure.
3. Create worker pools per domain (I/O, embeddings, retrieval, tools, maintenance).
4. Implement internal job scheduler.
5. Add cancellation/timeout with `context.Context` + `errgroup`.
6. Measure queue depths and per-stage latency.
7. Run `go test -race` on all core packages.

## Agents Used

- `go-concurrency-perf`
- `architect-go-systems`

## Outputs

- `internal/core/` package (eventbus, queue, workers, scheduler)
- Race-free test suite
- Latency/throughput metrics
