---
name: tailwind-patterns
description: Tailwind CSS v4 principles. CSS-first configuration, container
  queries, modern patterns, design token architecture.
allowed-tools: Read Write Edit Glob Grep
---
# Tailwind CSS Patterns (v4 - 2025)

> Modern utility-first CSS with CSS-native configuration.

---

## 1. Tailwind v4 Architecture

### What Changed from v3

### v4 Core Concepts

---

## 2. CSS-Based Configuration

### Theme Definition

```
@theme {
  /* Colors - use semantic names */
  --color-primary: oklch(0.7 0.15 250);
  --color-surface: oklch(0.98 0 0);
  --color-surface-dark: oklch(0.15 0 0);

  /* Spacing scale */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 2rem;

  /* Typography */
  --font-sans: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;
}
```

### When to Extend vs Override

---

## 3. Container Queries (v4 Native)

### Breakpoint vs Container

### Container Query Usage

### When to Use

---

## 4. Responsive Design

### Breakpoint System

### Mobile-First Principle

1. Write mobile styles first (no prefix)
2. Add larger screen overrides with prefixes
3. Example: `w-full md:w-1/2 lg:w-1/3`

---

## 5. Dark Mode

### Configuration Strategies

### Dark Mode Pattern

---

## 6. Modern Layout Patterns

### Flexbox Patterns

### Grid Patterns

> **Note:** Prefer asymmetric/Bento layouts over symmetric 3-column grids.

---

## 7. Modern Color System

### OKLCH vs RGB/HSL

### Color Token Architecture

---

## 8. Typography System

### Font Stack Pattern

### Type Scale

---

## 9. Animation & Transitions

### Built-in Animations

### Transition Patterns

---

## 10. Component Extraction

### When to Extract

### Extraction Methods

---

## 11. Anti-Patterns

---

## 12. Performance Principles

---

> **Remember:** Tailwind v4 is CSS-first. Embrace CSS variables, container queries, and native features. The config file is now optional.