---
name: server-management
description: Server management principles and decision-making. Process
  management, monitoring strategy, and scaling decisions. Teaches thinking, not
  commands.
allowed-tools: Read Write Edit Glob Grep Bash
---
# Server Management

> Server management principles for production operations.
> **Learn to THINK, not memorize commands.**

---

## 1. Process Management Principles

### Tool Selection

### Process Management Goals

---

## 2. Monitoring Principles

### What to Monitor

### Alert Severity Strategy

### Monitoring Tool Selection

---

## 3. Log Management Principles

### Log Strategy

### Log Principles

1. **Rotate logs** to prevent disk fill
2. **Structured logging** (JSON) for parsing
3. **Appropriate levels** (error/warn/info/debug)
4. **No sensitive data** in logs

---

## 4. Scaling Decisions

### When to Scale

### Scaling Strategy

---

## 5. Health Check Principles

### What Constitutes Healthy

### Health Check Implementation

- Simple: Just return 200
- Deep: Check all dependencies
- Choose based on load balancer needs

---

## 6. Security Principles

---

## 7. Troubleshooting Priority

When something's wrong:

1. **Check if running** (process status)
2. **Check logs** (error messages)
3. **Check resources** (disk, memory, CPU)
4. **Check network** (ports, DNS)
5. **Check dependencies** (database, APIs)

---

## 8. Anti-Patterns

---

> **Remember:** A well-managed server is boring. That's the goal.