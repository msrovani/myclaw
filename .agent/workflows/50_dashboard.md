---
description: "Phase 7 — Build operational dashboard with templ + HTMX + ECharts."
---

# 50 — Dashboard

## Steps

1. Build admin API endpoints (Go `net/http`).
2. Create templ templates for dashboard layout.
3. Implement HTMX partial updates for real-time data.
4. Add ECharts for throughput, tokens, memory, router charts.
5. Build panels: system overview, tokens/cost, memory, LLM router, skills, admin.
6. Implement admin controls (feature flags, retention, tuning, maintenance).
7. Add responsive CSS + dark mode.
8. Browser testing.

## Agents Used

- `dashboard-ops-ui`
- `architect-go-systems`

## Outputs

- `internal/dashboard/` package
- `web/` templates and assets
- Admin API handlers
- Browser test results
