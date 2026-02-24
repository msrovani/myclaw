---
name: planning-with-files
description: Implements Manus-style file-based planning for complex tasks.
  Creates task_plan.md, findings.md, and progress.md. Use when starting complex
  multi-step tasks, research projects, or any task requiring >5 tool calls.
allowed-tools: Read Write Edit Bash Glob Grep WebFetch WebSearch
metadata:
  version: 2.1.2
  user-invocable: "true"
  hooks: "{\"SessionStart\":[{\"hooks\":[{\"type\":\"command\",\"command\":\"echo
    '[planning-with-files] Ready. Auto-activates for complex tasks, or invoke
    manually with
    /planning-with-files'\"}]}],\"PreToolUse\":[{\"matcher\":\"Write|Edit|Bash\
    \",\"hooks\":[{\"type\":\"command\",\"command\":\"cat task_plan.md
    2>/dev/null | head -30 ||
    true\"}]}],\"PostToolUse\":[{\"matcher\":\"Write|Edit\",\"hooks\":[{\"type\
    \":\"command\",\"command\":\"echo '[planning-with-files] File updated. If
    this completes a phase, update task_plan.md
    status.'\"}]}],\"Stop\":[{\"hooks\":[{\"type\":\"command\",\"command\":\"${\
    CLAUDE_PLUGIN_ROOT}/scripts/check-complete.sh\"}]}]}"
---
# Planning with Files

Work like Manus: Use persistent markdown files as your "working memory on disk."

## Important: Where Files Go

When using this skill:

- **Templates** are stored in the skill directory at `${CLAUDE_PLUGIN_ROOT}/templates/`
- **Your planning files** (`task_plan.md`, `findings.md`, `progress.md`) should be created in **your project directory** — the folder where you're working

This ensures your planning files live alongside your code, not buried in the skill installation folder.

## Quick Start

Before ANY complex task:

1. **Create `task_plan.md`** in your project — Use [templates/task_plan.md](templates/task_plan.md) as reference
2. **Create `findings.md`** in your project — Use [templates/findings.md](templates/findings.md) as reference
3. **Create `progress.md`** in your project — Use [templates/progress.md](templates/progress.md) as reference
4. **Re-read plan before decisions** — Refreshes goals in attention window
5. **Update after each phase** — Mark complete, log errors

> **Note:** All three planning files should be created in your current working directory (your project root), not in the skill's installation folder.

## The Core Pattern

```
Context Window = RAM (volatile, limited)
Filesystem = Disk (persistent, unlimited)

→ Anything important gets written to disk.
```

## File Purposes

## Critical Rules

### 1. Create Plan First
Never start a complex task without `task_plan.md`. Non-negotiable.

### 2. The 2-Action Rule
> "After every 2 view/browser/search operations, IMMEDIATELY save key findings to text files."

This prevents visual/multimodal information from being lost.

### 3. Read Before Decide
Before major decisions, read the plan file. This keeps goals in your attention window.

### 4. Update After Act
After completing any phase:
- Mark phase status: `in_progress` → `complete`
- Log any errors encountered
- Note files created/modified

### 5. Log ALL Errors
Every error goes in the plan file. This builds knowledge and prevents repetition.

```markdown
## Errors Encountered
```

### 6. Never Repeat Failures
```
if action_failed:
    next_action != same_action
```
Track what you tried. Mutate the approach.

## The 3-Strike Error Protocol

```
ATTEMPT 1: Diagnose & Fix
  → Read error carefully
  → Identify root cause
  → Apply targeted fix

ATTEMPT 2: Alternative Approach
  → Same error? Try different method
  → Different tool? Different library?
  → NEVER repeat exact same failing action

ATTEMPT 3: Broader Rethink
  → Question assumptions
  → Search for solutions
  → Consider updating the plan

AFTER 3 FAILURES: Escalate to User
  → Explain what you tried
  → Share the specific error
  → Ask for guidance
```

## Read vs Write Decision Matrix

## The 5-Question Reboot Test

If you can answer these, your context management is solid:

## When to Use This Pattern

**Use for:**
- Multi-step tasks (3+ steps)
- Research tasks
- Building/creating projects
- Tasks spanning many tool calls
- Anything requiring organization

**Skip for:**
- Simple questions
- Single-file edits
- Quick lookups

## Templates

Copy these templates to start:

- [templates/task_plan.md](templates/task_plan.md) — Phase tracking
- [templates/findings.md](templates/findings.md) — Research storage
- [templates/progress.md](templates/progress.md) — Session logging

## Scripts

Helper scripts for automation:

- `scripts/init-session.sh` — Initialize all planning files
- `scripts/check-complete.sh` — Verify all phases complete

## Advanced Topics

- **Manus Principles:** See [reference.md](reference.md)
- **Real Examples:** See [examples.md](examples.md)

## Anti-Patterns