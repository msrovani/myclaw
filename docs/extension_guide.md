# XXXCLAW Extension Guide

XXXCLAW is designed to be easily extensible. This guide covers the two main expansion points: adding new LLM Providers and new Core Skills.

## 1. Adding a new LLM Provider

Providers are located in `internal/providers/`. Expanding the LLM ecosystem is simply a matter of satisfying the `Provider` interface:

```go
type Provider interface {
 ID() string
 Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
 Embed(ctx context.Context, text string) ([]float32, error)
}
```

### Steps

1. Create a `NewMyProvider(apiUrl string) *MyProvider` in the `providers` package.
2. In `cmd/xxxclaw/main.go`, instantiate your provider and push it into the initialization slice of `router.NewRouter()`.
3. The routing engine will immediately begin utilizing your model under the economy tracking metrics, routing jobs to it when policy flags align.

## 2. Adding a Core Go Skill

XXXCLAW registers dynamically created `.agent` file integrations under the hood, but sometimes raw performance or secure system access dictates writing a localized Go tool.

Skills live in `internal/skills/` and must match the `Skill` interface:

```go
type Skill interface {
 ID() string
 Description() string
 Execute(ctx context.Context, req Request) (Response, error)
}
```

### Steps

1. Create a new struct embedding your specific external libraries or API drivers.
2. Ensure you extract the TenantContext bounds:

```go
 tc, err := core.TenantFromContext(ctx)
 if err != nil {
  return Response{Error: "unauthorized"}, err
 }
```

3. Process the payload via `ParsePayload(req, &myStruct{})`.
2. Register the Skill in `cmd/xxxclaw/main.go` via `skillRegistry.Register(mySkill)`.

By leveraging these points, XXXCLAW acts as a foundational operating system for Autonomous AI Workloads.
