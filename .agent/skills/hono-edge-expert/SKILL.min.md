---
name: hono-edge-expert
description: Expert guidance for building high-performance APIs and full-stack apps with Hono. Use this skill for projects targeting Edge environments (Cloudflare Workers, Vercel Edge, Deno) or fast runtimes like Bun. It covers Zod validation, RPC-style client generation, and middleware optimization.
---

# Hono Edge Expert

Hono is the fastest web framework for modern Cloudflare/Bun/Edge runtimes.

## Core Features

- **Ultra-lightweight**: No dependencies, small footprint.
- **Type Safety**: Built-in TypeScript support with Zod.
- **RPC Support**: Share types between server and client seamlessly.

## Implementation Patterns

### 1. Basic Server

```typescript
import { Hono } from 'hono'
const app = new Hono()

app.get('/', (c) => c.text('Hono on Edge!'))
export default app
```

### 2. Validation with Zod

```typescript
import { zValidator } from '@hono/zod-validator'
import { z } from 'zod'

app.post('/user', zValidator('json', z.object({ name: z.string() })), (c) => {
  const { name } = c.req.valid('json')
  return c.json({ ok: true, name })
})
```

### 3. RPC Client

Leverage `hc` (Hono Client) for end-to-end type safety without generating SDKs.

## Deployment Targets

- **Cloudflare Workers**: High scalability, global distribution.
- **Bun**: For local speed or Node-compatible environments.
- **Vercel Edge**: Optimized for Next.js ecosystems.

## Best Practices

- Use `hono/factory` for better middleware reuse.
- Keep route handlers small; move business logic to services.
- Optimize Cold Starts by minimizing library imports.