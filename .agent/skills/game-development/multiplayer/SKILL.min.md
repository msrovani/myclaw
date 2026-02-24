---
name: multiplayer
description: Multiplayer game development principles. Architecture, networking, synchronization.
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# Multiplayer Game Development

> Networking architecture and synchronization principles.

---

## 1. Architecture Selection

### Decision Tree

```
What type of multiplayer?
│
├── Competitive / Real-time
│   └── Dedicated Server (authoritative)
│
├── Cooperative / Casual
│   └── Host-based (one player is server)
│
├── Turn-based
│   └── Client-server (simple)
│
└── Massive (MMO)
    └── Distributed servers
```

### Comparison

---

## 2. Synchronization Principles

### State vs Input

### Lag Compensation

---

## 3. Network Optimization

### Bandwidth Reduction

### Update Rates

---

## 4. Security Principles

### Server Authority

```
Client: "I hit the enemy"
Server: Validate → did projectile actually hit?
         → was player in valid state?
         → was timing possible?
```

### Anti-Cheat

---

## 5. Matchmaking

### Considerations

---

## 6. Anti-Patterns

---

> **Remember:** Never trust the client. The server is the source of truth.