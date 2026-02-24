---
name: tdd-workflow
description: Test-Driven Development workflow principles. RED-GREEN-REFACTOR cycle.
allowed-tools: Read Write Edit Glob Grep Bash
---
# TDD Workflow

> Write tests first, code second.

---

## 1. The TDD Cycle

```
🔴 RED → Write failing test
    ↓
🟢 GREEN → Write minimal code to pass
    ↓
🔵 REFACTOR → Improve code quality
    ↓
   Repeat...
```

---

## 2. The Three Laws of TDD

1. Write production code only to make a failing test pass
2. Write only enough test to demonstrate failure
3. Write only enough code to make the test pass

---

## 3. RED Phase Principles

### What to Write

### RED Phase Rules

- Test must fail first
- Test name describes expected behavior
- One assertion per test (ideally)

---

## 4. GREEN Phase Principles

### Minimum Code

### GREEN Phase Rules

- Don't write unneeded code
- Don't optimize yet
- Pass the test, nothing more

---

## 5. REFACTOR Phase Principles

### What to Improve

### REFACTOR Rules

- All tests must stay green
- Small incremental changes
- Commit after each refactor

---

## 6. AAA Pattern

Every test follows:

---

## 7. When to Use TDD

---

## 8. Test Prioritization

---

## 9. Anti-Patterns

---

## 10. AI-Augmented TDD

### Multi-Agent Pattern

---

> **Remember:** The test is the specification. If you can't write a test, you don't understand the requirement.