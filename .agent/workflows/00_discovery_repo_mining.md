---
description: "Phase 0 — Study reference repos, extract portable patterns, generate mining report."
---

# 00 — Discovery & Repo Mining

## Steps

1. List all reference repositories (Mem0, sqlite-vec, mem0-mcp, workspace `.agent`).
2. For each repo, extract:
   - Architecture patterns
   - Memory API contracts
   - Concurrency patterns
   - Skill/workflow conventions
   - Dashboard UX
   - Observability approach
   - License
3. Generate `docs/repo_mining.md` with comparative matrix.
4. Generate `docs/license_matrix.md`.
5. Identify patterns worth porting to Go — document rationale.

## Agents Used

- `explorer-agent`
- `architect-go-systems`

## Outputs

- `docs/repo_mining.md`
- `docs/license_matrix.md`
