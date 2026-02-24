---
name: powershell-windows
description: PowerShell Windows patterns. Critical pitfalls, operator syntax,
  error handling.
allowed-tools: Read Write Edit Glob Grep Bash
---
# PowerShell Windows Patterns

> Critical patterns and pitfalls for Windows PowerShell.

---

## 1. Operator Syntax Rules

### CRITICAL: Parentheses Required

**Rule:** Each cmdlet call MUST be in parentheses when using logical operators.

---

## 2. Unicode/Emoji Restriction

### CRITICAL: No Unicode in Scripts

**Rule:** Use ASCII characters only in PowerShell scripts.

---

## 3. Null Check Patterns

### Always Check Before Access

---

## 4. String Interpolation

### Complex Expressions

**Pattern:**
```
$value = $obj.prop.sub
Write-Output "Value: $value"
```

---

## 5. Error Handling

### ErrorActionPreference

### Try/Catch Pattern

- Don't return inside try block
- Use finally for cleanup
- Return after try/catch

---

## 6. File Paths

### Windows Path Rules

**Rule:** Use Join-Path for cross-platform safety.

---

## 7. Array Operations

### Correct Patterns

---

## 8. JSON Operations

### CRITICAL: Depth Parameter

**Rule:** Always specify `-Depth` for nested objects.

### File Operations

---

## 9. Common Errors

---

## 10. Script Template

```powershell
# Strict mode
Set-StrictMode -Version Latest
$ErrorActionPreference = "Continue"

# Paths
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Main
try {
    # Logic here
    Write-Output "[OK] Done"
    exit 0
}
catch {
    Write-Warning "Error: $_"
    exit 1
}
```

---

> **Remember:** PowerShell has unique syntax rules. Parentheses, ASCII-only, and null checks are non-negotiable.