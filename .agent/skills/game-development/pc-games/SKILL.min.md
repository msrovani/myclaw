---
name: pc-games
description: PC and console game development principles. Engine selection, platform features, optimization strategies.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# PC/Console Game Development

> Engine selection and platform-specific principles.

---

## 1. Engine Selection

### Decision Tree

```
What are you building?
│
├── 2D Game
│   ├── Open source important? → Godot
│   └── Large team/assets? → Unity
│
├── 3D Game
│   ├── AAA visual quality? → Unreal
│   ├── Cross-platform priority? → Unity
│   └── Indie/open source? → Godot 4
│
└── Specific Needs
    ├── DOTS performance? → Unity
    ├── Nanite/Lumen? → Unreal
    └── Lightweight? → Godot
```

### Comparison

---

## 2. Platform Features

### Steam Integration

### Console Requirements

---

## 3. Controller Support

### Input Abstraction

```
Map ACTIONS, not buttons:
- "confirm" → A (Xbox), Cross (PS), B (Nintendo)
- "cancel" → B (Xbox), Circle (PS), A (Nintendo)
```

### Haptic Feedback

---

## 4. Performance Optimization

### Profiling First

### Common Bottlenecks

---

## 5. Engine-Specific Principles

### Unity 6

- DOTS for performance-critical systems
- Burst compiler for hot paths
- Addressables for asset streaming

### Godot 4

- GDScript for rapid iteration
- C# for complex logic
- Signals for decoupling

### Unreal 5

- Blueprint for designers
- C++ for performance
- Nanite for high-poly environments
- Lumen for dynamic lighting

---

## 6. Anti-Patterns

---

> **Remember:** Engine is a tool. Master the principles, then adapt to any engine.