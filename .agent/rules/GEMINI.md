---
trigger: always_on
---

# GEMINI.md - Antigravity Kit (Lite)

> 🔴 **MANDATORY:** Read `.agent2/SYSTEM_LITE.md` and `.agent2/ARCHITECTURE_LITE.md` first.

## PROTOCOL
- Agent activated → Load skills from frontmatter → Read `.agent2/skills/<skill>/SKILL.min.md`.
- **Minified Context:** Prioritize `.min.md` versions to save tokens.
- **Selective Execution:** Only execute scripts via Python when strictly necessary.

## REQUEST CLASSIFIER (Concise)
1. **QUESTION**: Text only.
2. **SIMPLE CODE**: Single file edit.
3. **COMPLEX**: `{task-slug}.md` + Agent.
4. **SLASH**: Follow workflow.

## RULES
- **Language**: Translate internally, respond in user's language.
- **Code**: Clean Code skill mandatory. AAA Testing.
- **Socratic Gate**: 3 questions for build/complex tasks.
- **Final Checks**: Run `.agent2/scripts/checklist.py`.

## AGENT ROUTING
Identify domain → Open `.agent2/agents/{agent}.md` → Announce `🤖 Applying knowledge of @[agent]...`.
