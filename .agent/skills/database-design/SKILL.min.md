---
name: database-design
description: Database design principles and decision-making. Schema design,
  indexing strategy, ORM selection, serverless databases.
allowed-tools: Read Write Edit Glob Grep
---
# Database Design

> **Learn to THINK, not copy SQL patterns.**

## ⚠️ Core Principle

- ASK user for database preferences when unclear
- Choose database/ORM based on CONTEXT
- Don't default to PostgreSQL for everything

---

## Decision Checklist

Before designing schema:

- [ ] Asked user about database preference?
- [ ] Chosen database for THIS context?
- [ ] Considered deployment environment?
- [ ] Planned index strategy?
- [ ] Defined relationship types?

---

## Anti-Patterns

❌ Default to PostgreSQL for simple apps (SQLite may suffice)
❌ Skip indexing
❌ Use SELECT * in production
❌ Store JSON when structured data is better
❌ Ignore N+1 queries