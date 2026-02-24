---
name: threejs-shader-mastery
description: Advanced guidance for creating high-impact 3D experiences and custom GLSL shaders with Three.js. Use this skill for visual-heavy projects requiring custom post-processing, physics-based materials, and GPU-accelerated animations.
---

# Three.js & Shader Mastery

Transform standard 3D web scenes into cinematic experiences using custom GLSL and advanced Three.js techniques.

## Custom Shaders (GLSL)

### 1. Vertex Shader

Handles geometry and vertex positions. Use for ripple effects, flag waving, or mesh deformation.

### 2. Fragment Shader

Handles pixel color and lighting. Use for glass effects, gradients, procedural textures, and bloom.

```glsl
// Example Fragment Shader
varying vec2 vUv;
uniform float uTime;
void main() {
  gl_FragColor = vec4(vUv.x, sin(uTime), vUv.y, 1.0);
}
```

## Advanced Techniques

- **Post-Processing**: Using `EffectComposer` for gloom, depth-of-field, and retro-CRT effects.
- **InstancedMesh**: Drawing thousands of objects in a single draw call for high performance.
- **BufferGeometry**: Optimized data handling for large meshes.

## Performance Optimization

- **GPU Offloading**: Move animations from JS `requestAnimationFrame` to Shaders.
- **Texture Packing**: Combine multiple textures into a single sprite sheet.
- **Level of Detail (LOD)**: Reducing mesh complexity based on distance from camera.

## Tools

- **ShaderToy**: For prototyping GLSL code.
- **Three-Custom-Shader-Material**: For extending built-in materials.