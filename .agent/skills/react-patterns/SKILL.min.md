---
name: react-patterns
description: Modern React patterns and principles. Hooks, composition,
  performance, TypeScript best practices.
allowed-tools: Read Write Edit Glob Grep
---
# React Patterns

> Principles for building production-ready React applications.

---

## 1. Component Design Principles

### Component Types

### Design Rules

- One responsibility per component
- Props down, events up
- Composition over inheritance
- Prefer small, focused components

---

## 2. Hook Patterns

### When to Extract Hooks

### Hook Rules

- Hooks at top level only
- Same order every render
- Custom hooks start with "use"
- Clean up effects on unmount

---

## 3. State Management Selection

### State Placement

---

## 4. React 19 Patterns

### New Hooks

### Compiler Benefits

- Automatic memoization
- Less manual useMemo/useCallback
- Focus on pure components

---

## 5. Composition Patterns

### Compound Components

- Parent provides context
- Children consume context
- Flexible slot-based composition
- Example: Tabs, Accordion, Dropdown

### Render Props vs Hooks

---

## 6. Performance Principles

### When to Optimize

### Optimization Order

1. Check if actually slow
2. Profile with DevTools
3. Identify bottleneck
4. Apply targeted fix

---

## 7. Error Handling

### Error Boundary Usage

### Error Recovery

- Show fallback UI
- Log error
- Offer retry option
- Preserve user data

---

## 8. TypeScript Patterns

### Props Typing

### Common Types

---

## 9. Testing Principles

### Test Priorities

- User-visible behavior
- Edge cases
- Error states
- Accessibility

---

## 10. Anti-Patterns

---

> **Remember:** React is about composition. Build small, combine thoughtfully.