---
name: clean-code
description: Pragmatic coding standards - concise, direct, no over-engineering,
  no unnecessary comments
allowed-tools: Read Write Edit
metadata:
  version: "2"
  priority: CRITICAL
---
# Clean Code - Pragmatic AI Coding Standards

> **CRITICAL SKILL** - Be **concise, direct, and solution-focused**.

---

## Core Principles

---

## Naming Rules

> **Rule:** If you need a comment to explain a name, rename it.

---

## Function Rules

---

## Code Structure

---

## AI Coding Style

---

## Anti-Patterns (DON'T)

---

## 🔴 Before Editing ANY File (THINK FIRST!)

**Before changing a file, ask yourself:**

**Quick Check:**

```
File to edit: UserService.ts
└── Who imports this? → UserController.ts, AuthController.ts
└── Do they need changes too? → Check function signatures
```

> 🔴 **Rule:** Edit the file + all dependent files in the SAME task.
> 🔴 **Never leave broken imports or missing updates.**

---

## Summary

> **Remember: The user wants working code, not a programming lesson.**

---

## 🔴 Self-Check Before Completing (MANDATORY)

**Before saying "task complete", verify:**

> 🔴 **Rule:** If ANY check fails, fix it before completing.

---

## Verification Scripts (MANDATORY)

> 🔴 **CRITICAL:** Each agent runs ONLY their own skill's scripts after completing work.

### Agent → Script Mapping

> ❌ **WRONG:** `test-engineer` running `ux_audit.py`
> ✅ **CORRECT:** `frontend-specialist` running `ux_audit.py`

---

### 🔴 Script Output Handling (READ → SUMMARIZE → ASK)

**When running a validation script, you MUST:**

1. **Run the script** and capture ALL output
2. **Parse the output** - identify errors, warnings, and passes
3. **Summarize to user** in this format:

```markdown
## Script Results: [script_name.py]

### ❌ Errors Found (X items)
- [File:Line] Error description 1
- [File:Line] Error description 2

### ⚠️ Warnings (Y items)
- [File:Line] Warning description

### ✅ Passed (Z items)
- Check 1 passed
- Check 2 passed

**Should I fix the X errors?**
```

1. **Wait for user confirmation** before fixing
2. **After fixing** → Re-run script to confirm

> 🔴 **VIOLATION:** Running script and ignoring output = FAILED task.
> 🔴 **VIOLATION:** Auto-fixing without asking = Not allowed.
> 🔴 **Rule:** Always READ output → SUMMARIZE → ASK → then fix.