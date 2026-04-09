# AGENTS.md

## Purpose
This file defines concurrency rules, invariants, and implementation responsibilities.  
All agents must follow these rules when analyzing or modifying code.

---

# 🧠 Global Concurrency Principles

- Concurrency control exists to **protect invariants**, not to optimize performance.
- Always identify the **invariant first**, then choose the mechanism.
- Prefer **database-enforced constraints** for uniqueness and duplication prevention.
- Prefer **pessimistic locking** when:
    - a command rewrites a child collection
    - multiple writers could overwrite each other
- Prefer **optimistic locking** when:
    - concurrent writes should fail fast
    - conflicts are rare and retriable
- Do NOT rely on application-level checks for invariants unless backed by the database.
- Do NOT widen locking scope without explicitly stating the invariant being protected.

---

# 🏗️ Architectural Responsibilities

## Controllers
- Must NOT contain business logic
- Must NOT implement concurrency or locking logic

## Application Services
- Define transaction boundaries
- Coordinate use cases
- Decide which repository operations are executed

## Repositories
- Implement concurrency mechanisms
- Enforce locking strategy (optimistic / pessimistic)
- Translate persistence errors into domain/application errors

## Domain Models
- Must remain pure (no ORM, no DB concerns)
- Express invariants, not persistence details

---

# 🔒 Concurrency Invariants

## User
**Invariant:** one account per email  
**Mechanism:** UNIQUE constraint on `users.email`  
**Why:** registration is a create command keyed by a natural identifier.  
Database uniqueness is the source of truth. Application checks are advisory only.

---

## Organization
**Invariant:** an owner cannot create two organizations with the same name (case-insensitive)  
**Mechanism:** unique index on `(user_id, lower(name))`  
**Why:** duplicate prevention belongs in the schema; application-level checks are not authoritative.

---

## League
**Invariant:** an organization cannot create two leagues with the same name (case-insensitive)  
**Mechanism:** unique index on `(organization_id, lower(name))`  
**Why:** enforced atomically at the database level.

---

## League Membership Set
**Invariant:**
- membership add/remove must not lose unrelated updates
- a team can appear only once per league

**Mechanism:**
- `SELECT ... FOR UPDATE` on the parent league row for ALL membership mutations
- unique index on `(league_id, team_id)`

**Why:**
membership commands rewrite a child collection.  
Parent-row locking serializes conflicting writers and prevents lost updates.

---

## Team
**Invariant:** team creation must atomically create:
- team row
- organization membership
- user-role assignments

**Mechanism:** single transaction across all writes

**Why:** consistency spans multiple tables; atomicity is required.

---

## Season
**Invariant:** only one active season per league  
(active = pending, planned, in_progress, paused)

**Mechanism:** partial unique index on `seasons(league_id)` for active states

**Why:** cross-aggregate invariant must be enforced centrally by the database.

---

## Season State & Match Updates
**Invariant:** state transitions and score updates must not overwrite each other

**Mechanism:** optimistic locking using `seasons.version` in update predicates

**Why:**  
Season is the aggregate root.  
Concurrent updates should fail explicitly, not be serialized silently.

---

# ⚠️ Conflict & Error Handling

- UNIQUE constraint violations → return domain-level conflict errors
- Optimistic lock failures → return concurrency conflict (no silent overwrite)
- Do NOT auto-retry optimistic failures unless explicitly required
- Pessimistic operations must:
    - occur inside a transaction
    - include read + mutate + write in the same transaction

---

# 🧪 Testing Requirements

For any concurrency-related change:
- Add tests for:
    - concurrent writes
    - duplicate creation attempts
    - invariant violations
- Tests must prove the invariant is protected, not just that code executes

---

# 🤖 Agent Behavior Rules

Agents MUST:

- Identify the invariant before modifying code
- Explain the concurrency risk before proposing changes
- Choose the **minimal mechanism** required to protect the invariant
- Avoid introducing broader locks than necessary
- Keep changes **scoped and minimal**
- Update tests when modifying concurrency logic

Agents MUST NOT:

- Replace database constraints with application checks
- Add locking without explaining the invariant it protects
- Spread concurrency logic across controllers/services arbitrarily
- Modify unrelated modules

---

# 🔄 Decision Framework (MANDATORY)

For any concurrency decision, agents must follow:

1. What invariant is being protected?
2. What race condition can violate it?
3. What is the minimal mechanism to prevent it?
    - UNIQUE constraint
    - Transaction
    - Pessimistic lock
    - Optimistic lock
    - Idempotency
4. Why is this better than alternatives?

---

# 📌 Summary

- Invariants drive design
- Database enforces truth
- Locks are used sparingly and intentionally
- Optimistic for state, pessimistic for collections
- Always prove correctness with tests
