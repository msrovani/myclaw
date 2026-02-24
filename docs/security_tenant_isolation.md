# Security & Tenant Isolation

This file defines the security posture for XXXCLAW tenants.

## Deny by Default

- No API endpoint, database function, skill execution, or memory retrieval will be permitted without explicit extraction and validation of `uid` and `workspace_id`.

## Path Traversal Protection

- DB files are dynamically constructed using `{uid}` and `{workspace_id}`.
- All IDs must be strictly validated (UUIDs, or alphanumeric bounds) to prevent directory traversal attacks (e.g. `workspace_id` = `../../../etc`).

## Defense in Depth (Dual Layer Isolation)

1. **File System Level**: Because each workspace corresponds to a unique `.db` file on disk, one tenant's SQL injection cannot leak into another tenant's file. Connections are opened specifically to the tenant's DB file.
2. **Schema Level**: All SQLite tables, even when isolated in their own physical DB file, include `uid` and `workspace_id` columns. All queries must include `WHERE uid=? AND workspace_id=?`.

## Logging and Dashboard

- `slog` output must strictly redact actual content.
- Observability tags/metrics must append `uid` or aggregate safely without leaking.
- Dashboards are strictly tenant-aware, isolating charts and statistics by verified headers.
