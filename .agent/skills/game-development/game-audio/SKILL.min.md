---
name: game-audio
description: Game audio principles. Sound design, music integration, adaptive audio systems.
allowed-tools: Read, Glob, Grep
---

# Game Audio Principles

> Sound design and music integration for immersive game experiences.

---

## 1. Audio Category System

### Category Definitions

### Priority Hierarchy

```
When sounds compete for channels:

1. Voice (highest - always audible)
2. Player SFX (feedback critical)
3. Enemy SFX (gameplay important)
4. Music (mood, but duckable)
5. Ambient (lowest - can drop)
```

---

## 2. Sound Design Decisions

### SFX Creation Approach

### Layering Structure

---

## 3. Music Integration

### Music State System

```
Game State → Music Response
│
├── Menu → Calm, loopable theme
├── Exploration → Ambient, atmospheric
├── Combat detected → Transition to tension
├── Combat engaged → Full battle music
├── Victory → Stinger + calm transition
├── Defeat → Somber stinger
└── Boss → Unique, multi-phase track
```

### Transition Techniques

---

## 4. Adaptive Audio Decisions

### Intensity Parameters

### Vertical vs Horizontal

---

## 5. 3D Audio Decisions

### Spatialization

### Distance Behavior

---

## 6. Platform Considerations

### Format Selection

### Memory Budget

---

## 7. Mix Hierarchy

### Volume Balance Reference

### Ducking Rules

---

## 8. Anti-Patterns

---

> **Remember:** 50% of the game experience is audio. A muted game loses half its soul.