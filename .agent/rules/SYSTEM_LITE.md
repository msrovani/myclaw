# SYSTEM_LITE.md - Antigravity Kit (Token Optimized)

## 🎯 CORE PROTOCOL
1. **Identify Agent**: Select from 20 specialists based on domain.
2. **Announce**: `🤖 **Applying knowledge of @[agent]...**`
3. **Load Skill**: Read `.agent2/skills/<skill>/SKILL.min.md` (ONLY what is requested).
4. **Socratic Gate**: Ask 3 strategic questions for new features/complex tasks.

## 📥 CLASSIFIER
- **QUESTION**: TIER 0 (Text only)
- **CODE/DESIGN**: TIER 1 (Agent + {task-slug}.md required for complex tasks)
- **SLASH CMD**: Execute specific workflow.

## 🧹 CLEAN CODE
- Concise, self-documenting, no over-engineering.
- Mandatory Testing (Unit > Int > E2E).
- Verify imports and dependencies before editing.

## 🏁 VALIDATION
- Run `python .agent2/scripts/checklist.py .` before completion.
- P0: Security/Lint > P1: Schema/Tests > P2: UX/SEO.

## 🎭 MODES
- **plan**: 4-Phase (Analysis, Planning, Solutioning, Implementation). No code before Phase 4.
- **edit**: Direct execution. Use {task-slug}.md for multi-file changes.
