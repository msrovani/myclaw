# Testing Tenant Isolation

This document describes the mandatory isolation test suite required before any release of XXXCLAW.

## Mandatory Test Cases

### Test A: Cross-User Memory Isolation

- **Setup**: Create User A and User B. Create Memory 1 for User A.
- **Action**: User B attempts to query Memory 1 (by ID or vector search).
- **Expectation**: Query returns 0 results or "Not Found".

### Test B: Cross-Workspace Isolation (Same User)

- **Setup**: Create User A with Workspace 1 and Workspace 2. Create Memory 1 in Workspace 1.
- **Action**: User A queries memory while authenticated under Workspace 2.
- **Expectation**: Query returns 0 results.

### Test C: Agent Scope Isolation

- **Setup**: User A, Workspace 1. Agent X and Agent Y. Agent X creates an `agent_scoped` memory.
- **Action**: Agent Y queries memory.
- **Expectation**: Query returns 0 results unless an explicit ACL policy grants Agent Y access.

### Test D: Vector Retrieval Isolation

- **Setup**: Inject embeddings into Workspace A's sqlite-vec index and Workspace B's index.
- **Action**: Search vectors in Workspace B using a term present only in Workspace A.
- **Expectation**: No semantic matches returned from Workspace A's index.

### Test E: Semantic Cache Isolation

- **Setup**: User A asks "What is the capital of France?" -> Cached in Workspace A.
- **Action**: User B asks the identical question.
- **Expectation**: User B query results in a cache MISS. User B should not receive User A's cached response.

### Test F: Global Metadata Isolation

- **Setup**: Review `data/control/control.db`.
- **Action**: Scan for PII, memory content, or private strings.
- **Expectation**: Zero private data. Only global system state (auth mappings, provider endpoints, etc).
