---
description: "Phase 9 — Security hardening, release engineering, final docs."
---

# 70 — Hardening & Release

## Steps

1. Security audit: multi-tenant isolation, input validation, SQL injection prevention.
2. Implement audit trail for sensitive operations.
3. Add log redaction for PII/secrets.
4. Verify feature flags (cloud disabled, dangerous tools deny-by-default).
5. Create Dockerfile (multi-stage build).
6. Create release Makefile targets.
7. Write deployment guide.
8. Generate changelog.
9. Create extension guide (how to add skill/tool/provider/memory policy).

## Agents Used

- `qa-security-release`
- `architect-go-systems`

## Outputs

- Security audit report
- Dockerfile
- Release Makefile targets
- `docs/deployment.md`
- `docs/extension-guide.md`
- `CHANGELOG.md`
