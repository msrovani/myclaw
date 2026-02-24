# XXXCLAW — License Matrix

## Direct Dependencies

| Dependency | License | Usage | Attribution Required |
|---|---|---|---|
| `github.com/asg017/sqlite-vec/bindings/go` | MIT + Apache-2.0 | Vector search extension | ✅ NOTICE file |
| `github.com/mattn/go-sqlite3` | MIT | SQLite3 Go driver (CGo) | ✅ |
| `github.com/a]h/templ` | MIT | HTML templating | ✅ |

## Reference Repos (Pattern Inspiration)

| Repository | License | Usage | Attribution |
|---|---|---|---|
| `mem0ai/mem0` | Apache-2.0 | API contract inspiration (not code) | ✅ NOTICE: API patterns |
| `mem0ai/mem0-mcp` | Apache-2.0 | MCP tool schema inspiration | ✅ NOTICE: MCP schema |
| `asg017/sqlite-vec` | MIT + Apache-2.0 | Direct dependency (Go bindings) | ✅ In go.mod + NOTICE |

## License Compliance Notes

1. **Apache-2.0** (Mem0): Allows use, modification, distribution. Must include NOTICE file if redistributing.
2. **MIT** (sqlite-vec, go-sqlite3, templ): Permissive. Include copyright notice.
3. **No code is copied** from reference repos. Only API contracts and architectural patterns are ported.
4. A `NOTICE` file will be maintained at project root crediting pattern inspirations.

## NOTICE File (to create at project root)

```
XXXCLAW
Copyright 2026 [Your Name/Org]

This project includes patterns inspired by:
- Mem0 (https://github.com/mem0ai/mem0) — Apache-2.0
  Memory API contract design (add/search/get/update/delete)
- sqlite-vec (https://github.com/asg017/sqlite-vec) — MIT + Apache-2.0
  Vector search integration via Go bindings
```
