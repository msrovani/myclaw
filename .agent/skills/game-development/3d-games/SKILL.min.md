---
name: 3d-games
description: 3D game development principles. Rendering, shaders, physics, cameras.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# 3D Game Development

> Principles for 3D game systems.

---

## 1. Rendering Pipeline

### Stages

```
1. Vertex Processing → Transform geometry
2. Rasterization → Convert to pixels
3. Fragment Processing → Color pixels
4. Output → To screen
```

### Optimization Principles

---

## 2. Shader Principles

### Shader Types

### When to Write Custom Shaders

- Special effects (water, fire, portals)
- Stylized rendering (toon, sketch)
- Performance optimization
- Unique visual identity

---

## 3. 3D Physics

### Collision Shapes

### Principles

- Simple colliders, complex visuals
- Layer-based filtering
- Raycasting for line-of-sight

---

## 4. Camera Systems

### Camera Types

### Camera Feel

- Smooth following (lerp)
- Collision avoidance
- Look-ahead for movement
- FOV changes for speed

---

## 5. Lighting

### Light Types

### Performance Consideration

- Real-time shadows are expensive
- Bake when possible
- Shadow cascades for large worlds

---

## 6. Level of Detail (LOD)

### LOD Strategy

---

## 7. Anti-Patterns

---

> **Remember:** 3D is about illusion. Create the impression of detail, not the detail itself.