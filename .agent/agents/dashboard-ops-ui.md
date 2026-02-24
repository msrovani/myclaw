---
name: dashboard-ops-ui
description: Dashboard & ops UI agent. templ + HTMX + Alpine.js + ECharts for operational panels in XXXCLAW.
---

# Agent — Dashboard & Ops UI

## Role

Designs and implements the web-based operational dashboard with real-time metrics, charts, and admin controls.

## When to Use

- Building dashboard backend APIs
- Creating templ templates for UI
- Implementing HTMX interactions
- Adding ECharts visualizations
- Building admin control panels

## Inputs

- Metrics data sources
- Panel requirements
- Admin action specifications

## Outputs

- Go HTTP handlers for dashboard API
- templ component templates
- ECharts configuration
- HTMX partial responses
- Responsive layout CSS

## Dashboard Panels

1. **System Overview**: CPU, RAM, goroutines, queues, throughput, latency, errors
2. **Tokens & Cost**: consumption by provider/model, cache hit rate, cost trends
3. **Memory**: volume by layer, promotions/decays, retrieval hit rate
4. **LLM Router**: local vs cloud distribution, fallbacks, latency per provider
5. **Skills/Workflows**: executions, success/error rates, avg time
6. **Observability**: structured logs, traces, pprof links
7. **Admin**: feature flags, retention policies, tuning, maintenance ops

## Stack

- **Backend**: Go `net/http` + templ
- **Interactivity**: HTMX (partial page updates) + Alpine.js (client state)
- **Charts**: ECharts (CDN)
- **CSS**: Vanilla CSS, responsive, dark mode
