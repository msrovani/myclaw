---
name: red-team-tactics
description: Red team tactics principles based on MITRE ATT&CK. Attack phases,
  detection evasion, reporting.
allowed-tools: Read Glob Grep
---
# Red Team Tactics

> Adversary simulation principles based on MITRE ATT&CK framework.

---

## 1. MITRE ATT&CK Phases

### Attack Lifecycle

```
RECONNAISSANCE → INITIAL ACCESS → EXECUTION → PERSISTENCE
       ↓              ↓              ↓            ↓
   PRIVILEGE ESC → DEFENSE EVASION → CRED ACCESS → DISCOVERY
       ↓              ↓              ↓            ↓
LATERAL MOVEMENT → COLLECTION → C2 → EXFILTRATION → IMPACT
```

### Phase Objectives

---

## 2. Reconnaissance Principles

### Passive vs Active

### Information Targets

---

## 3. Initial Access Vectors

### Selection Criteria

---

## 4. Privilege Escalation Principles

### Windows Targets

### Linux Targets

---

## 5. Defense Evasion Principles

### Key Techniques

### Operational Security

- Work during business hours
- Mimic legitimate traffic patterns
- Use encrypted channels
- Blend with normal behavior

---

## 6. Lateral Movement Principles

### Credential Types

### Movement Paths

- Admin shares
- Remote services (RDP, SSH, WinRM)
- Exploitation of internal services

---

## 7. Active Directory Attacks

### Attack Categories

---

## 8. Reporting Principles

### Attack Narrative

Document the full attack chain:
1. How initial access was gained
2. What techniques were used
3. What objectives were achieved
4. Where detection failed

### Detection Gaps

For each successful technique:
- What should have detected it?
- Why didn't detection work?
- How to improve detection

---

## 9. Ethical Boundaries

### Always

- Stay within scope
- Minimize impact
- Report immediately if real threat found
- Document all actions

### Never

- Destroy production data
- Cause denial of service (unless scoped)
- Access beyond proof of concept
- Retain sensitive data

---

## 10. Anti-Patterns

---

> **Remember:** Red team simulates attackers to improve defenses, not to cause harm.