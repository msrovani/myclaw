---
name: app-builder
description: Main application building orchestrator. Creates full-stack
  applications from natural language requests. Determines project type, selects
  tech stack, coordinates agents.
allowed-tools: Read Write Edit Glob Grep Bash Agent
---
# App Builder - Application Building Orchestrator

> Analyzes user's requests, determines tech stack, plans structure, and coordinates agents.

## 📦 Templates (13)

Quick-start scaffolding for new projects. **Read the matching template only!**

---

## 🔗 Related Agents

---

## Usage Example

```
User: "Make an Instagram clone with photo sharing and likes"

App Builder Process:
1. Project type: Social Media App
2. Tech stack: Next.js + Prisma + Cloudinary + Clerk
3. Create plan:
   ├─ Database schema (users, posts, likes, follows)
   ├─ API routes (12 endpoints)
   ├─ Pages (feed, profile, upload)
   └─ Components (PostCard, Feed, LikeButton)
4. Coordinate agents
5. Report progress
6. Start preview
```