---
name: game-development
description: Game development orchestrator. Routes to platform-specific skills
  based on project needs.
allowed-tools: Read Write Edit Glob Grep Bash
---
# Game Development

> **Orchestrator skill** that provides core principles and routes to specialized sub-skills.

---

## When to Use This Skill

You are working on a game development project. This skill teaches the PRINCIPLES of game development and directs you to the right sub-skill based on context.

---

## Sub-Skill Routing

### Platform Selection

### Dimension Selection

### Specialty Areas

---

## Core Principles (All Platforms)

### 1. The Game Loop

Every game, regardless of platform, follows this pattern:

```
INPUT  → Read player actions
UPDATE → Process game logic (fixed timestep)
RENDER → Draw the frame (interpolated)
```

**Fixed Timestep Rule:**
- Physics/logic: Fixed rate (e.g., 50Hz)
- Rendering: As fast as possible
- Interpolate between states for smooth visuals

---

### 2. Pattern Selection Matrix

**Decision Rule:** Start with State Machine. Add ECS only when performance demands.

---

### 3. Input Abstraction

Abstract input into ACTIONS, not raw keys:

```
"jump"  → Space, Gamepad A, Touch tap
"move"  → WASD, Left stick, Virtual joystick
```

**Why:** Enables multi-platform, rebindable controls.

---

### 4. Performance Budget (60 FPS = 16.67ms)

**Optimization Priority:**
1. Algorithm (O(n²) → O(n log n))
2. Batching (reduce draw calls)
3. Pooling (avoid GC spikes)
4. LOD (detail by distance)
5. Culling (skip invisible)

---

### 5. AI Selection by Complexity

---

### 6. Collision Strategy

---

## Anti-Patterns (Universal)

---

## Routing Examples

### Example 1: "I want to make a browser-based 2D platformer"
→ Start with `game-development/web-games` for framework selection
→ Then `game-development/2d-games` for sprite/tilemap patterns
→ Reference `game-development/game-design` for level design

### Example 2: "Mobile puzzle game for iOS and Android"
→ Start with `game-development/mobile-games` for touch input and stores
→ Use `game-development/game-design` for puzzle balancing

### Example 3: "Multiplayer VR shooter"
→ `game-development/vr-ar` for comfort and immersion
→ `game-development/3d-games` for rendering
→ `game-development/multiplayer` for networking

---

> **Remember:** Great games come from iteration, not perfection. Prototype fast, then polish.