# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - XXXCLAW Release

### Added

- **Core Runtime**: Go concurrent worker pools, internal bounded queue event bus.
- **Persistent Memory**: SQLite-backed layered memory (Short-term/Long-term), strictly isolated via TenantContexts (`uid`, `workspace_id`).
- **Hybrid Search**: Fast pure-Go vector cosine distance coupled with FTS5 search scoring using Reciprocal Rank Fusion (RRF).
- **LLM Routing**: Advanced provider interface with token tracking and budget enforcement isolating workloads per tenant.
- **Semantic Cache**: Reduces redundant token spend by serving identical prompts directly from memory.
- **Skills Registry**: Go interfaces for agent tools with panic-recovery, limits, audit logging, and dynamic `.agent` file integrations.
- **Dashboard**: Minimalist high-performance interface built with templ, HTMX, Alpine.js, and Echarts for operational oversight.

### Changed

- Shifted initial CGo vector roadmap to pure Go implementation to guarantee robust local Windows execution environments.
- Enforced strict multitenancy deep into the database manager, physically segregating tenant workspaces into individual DB files.
