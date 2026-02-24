---
name: game-art
description: Game art principles. Visual style selection, asset pipeline, animation workflow.
allowed-tools: Read, Glob, Grep
---

# Game Art Principles

> Visual design thinking for games - style selection, asset pipelines, and art direction.

---

## 1. Art Style Selection

### Decision Tree

```
What feeling should the game evoke?
│
├── Nostalgic / Retro
│   ├── Limited palette? → Pixel Art
│   └── Hand-drawn feel? → Vector / Flash style
│
├── Realistic / Immersive
│   ├── High budget? → PBR 3D
│   └── Stylized realism? → Hand-painted textures
│
├── Approachable / Casual
│   ├── Clean shapes? → Flat / Minimalist
│   └── Soft feel? → Gradient / Soft shadows
│
└── Unique / Experimental
    └── Define custom style guide
```

### Style Comparison Matrix

---

## 2. Asset Pipeline Decisions

### 2D Pipeline

### 3D Pipeline

---

## 3. Color Theory Decisions

### Palette Selection

### Color Principles

- **Hierarchy:** Important elements should pop
- **Consistency:** Same object = same color family
- **Context:** Colors read differently on backgrounds
- **Accessibility:** Don't rely only on color

---

## 4. Animation Principles

### The 12 Principles (Applied to Games)

### Frame Count Guidelines

---

## 5. Resolution & Scale Decisions

### 2D Resolution by Platform

### Consistency Rule

Choose a base unit and stick to it:
- Pixel art: Work at 1x, scale up (never down)
- HD art: Define DPI, maintain ratio
- 3D: 1 unit = 1 meter (industry standard)

---

## 6. Asset Organization

### Naming Convention

```
[type]_[object]_[variant]_[state].[ext]

Examples:
spr_player_idle_01.png
tex_stone_wall_normal.png
mesh_tree_oak_lod2.fbx
```

### Folder Structure Principle

```
assets/
├── characters/
│   ├── player/
│   └── enemies/
├── environment/
│   ├── props/
│   └── tiles/
├── ui/
├── effects/
└── audio/
```

---

## 7. Anti-Patterns

---

> **Remember:** Art serves gameplay. If it doesn't help the player, it's decoration.