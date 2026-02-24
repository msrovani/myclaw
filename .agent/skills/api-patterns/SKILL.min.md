---
name: api-patterns
description: API design principles and decision-making. REST vs GraphQL vs tRPC
  selection, response formats, versioning, pagination.
allowed-tools: Read Write Edit Glob Grep
---
# API Patterns

> API design principles and decision-making for 2025.
> **Learn to THINK, not copy fixed patterns.**

## ❌ Anti-Patterns

**DON'T:**
- Default to REST for everything
- Use verbs in REST endpoints (/getUsers)
- Return inconsistent response formats
- Expose internal errors to clients
- Skip rate limiting

**DO:**
- Choose API style based on context
- Ask about client requirements
- Document thoroughly
- Use appropriate status codes

---

## Script