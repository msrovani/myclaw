# Tenant Isolation Model & Architecture

## Core Principle

XXXCLAW is multi-agent and multi-tenant out of the box. Data isolation is strictly enforced.

- **"Secure by design"** and **"Deny by default"**.
- All memories, knowledge, artifacts, vector indexes, semantic caches, and histories belong **EXCLUSIVELY** to the `UID` and `WorkspaceID` that created them.

## 1. Context Contract (TenantContext)

Every read/write operation, database query, agent invocation, and workflow execution requires a valid `TenantContext`:

```go
type TenantContext struct {
    UID         string
    WorkspaceID string
    AgentID     string
    SessionID   string
    AuthClaims  []string
    RequestID   string
}
```

If this context is missing from the Go `context.Context`, the operation fails immediately. NO global operations are permitted.

## 2. Physical Database Isolation (DB-per-Workspace)

The primary mode of isolation is physical: **1 SQLite DB per Workspace**.

- Metadata/Control (Global without private tenant data): `data/control/control.db`
- Workspace DB: `data/tenants/{uid}/workspace/{workspace_id}/memory.db`
- Workspace Cache: `data/tenants/{uid}/workspace/{workspace_id}/cache.db`

Even within physical files, all tables enforce logical scoping logic (`uid` and `workspace_id` columns) for defense-in-depth and potential shared-DB deployments in Edge environments.

## 3. Scopes & ACL

Every memory or knowledge item has explicitly defined ACL properties:

- `owner_uid` (Required)
- `workspace_id` (Required)
- `visibility` (Default: `private`)
- `created_by_agent_id` (Optional)
- `access_policy`:
  - `owner_only`: Default
  - `workspace_agents`: Accessible to all agents in this workspace
  - `explicit_agent_list`: For precise graph boundaries

## 4. Multi-Agent Boundaries

Agents only run inside 1 Workspace at a time. They can maintain `agent_scoped` memories or read `workspace_scoped` memories, provided they belong to the identical `owner_uid` and `workspace_id`. Cross-user agent interaction is strictly prohibited at the memory level.
