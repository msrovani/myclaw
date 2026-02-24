---
name: nextjs-15-migration
description: Comprehensive guide for migrating to Next.js 15 and implementing its latest features. Use this skill when upgrading legacy Next.js apps or starting new projects that leverage React 19, the React Compiler, Partial Prerendering (PPR), and new caching behaviors.
---

# Next.js 15 Migration & Best Practices

Next.js 15 introduces React 19 support and major shifts in how applications are built and cached.

## Key Changes in V15

### 1. Caching is now 'Opt-In' by Default

- `fetch` requests are no longer cached by default (`cache: 'no-store'`).
- GET route handlers are no longer cached by default.

### 2. Async Request APIs

Headers and Cookies are now asynchronous.

```typescript
// Old
const cookieStore = cookies();
// New 
const cookieStore = await cookies();
```

### 3. React 19 & React Compiler

- Auto-memoization (no more `useMemo` / `useCallback` in most cases).
- Use `experimental: { reactCompiler: true }` in `next.config.js`.

### 4. Partial Prerendering (PPR)

Allows combining static and dynamic content in the same route.

- Set `experimental: { ppr: 'incremental' }`.

## Migration Workflow

1. **Update**: `npm install next@latest react@rc react-dom@rc`.
2. **Scan**: Identify synchronous `cookies()` and `headers()` calls.
3. **Refactor**: Update dynamic routing parameters (they are now Promises).
4. **Test**: Verify hydration and server component rendering with the new Compiler.

## Best Practices

- Prefer Server Actions for all data mutations.
- Use `Suspense` aggressively to enable PPR.
- Monitor bundle sizes as the Compiler changes tree-shaking behavior.