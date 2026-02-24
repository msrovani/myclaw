---
name: tauri-desktop-architecture
description: Architecture and development patterns for building lightweight, secure desktop applications using Tauri. Use this skill when building apps with a Rust backend and any frontend (React, Vue, Svelte). It covers Command handling, IPC security, and auto-update configurations.
---

# Tauri Desktop Architecture

Tauri is the modern, secure alternative to Electron, replacing Chromium with the native OS WebView.

## Architecture Concept

- **Frontend**: HTML/JS/CSS (React, Next.js, etc.) - Runs in a sandboxed WebView.
- **Backend (Core)**: Rust - Handles file system, system calls, and heavy processing.
- **IPC (Inter-Process Communication)**: JSON-based bridge between JS and Rust.

## Security Model

- **Isolator Pattern**: Prevents malicious scripts in the frontend from accessing the OS.
- **Capability System**: Explicitly allow specific system features (HTTP, FS, Shell) in `tauri.conf.json`.

## Implementation Patterns

### 1. Rust Commands

```rust
#[tauri::command]
fn greet(name: &str) -> String {
   format!("Hello, {}!", name)
}
```

### 2. Invoking from JS

```javascript
import { invoke } from '@tauri-apps/api/tauri'
const response = await invoke('greet', { name: 'User' })
```

## Key Benefits

- **Bundle Size**: ~3MB (vs 100MB+ for Electron).
- **Security**: Hardened by default.
- **Rust Power**: Native performance for intensive tasks.

## Best Practices

- Move long-running tasks to Rust side using `tokio`.
- Use a state management library in Rust for persistent data.
- Sign your binaries for production releases.