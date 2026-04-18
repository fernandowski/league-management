---
name: ddd-domain-modeling
description: Model business concepts using Domain-Driven Design principles. Classifies concepts into aggregate roots, entities, value objects, domain services, or application services and generates code following repository conventions.
---

## Use this skill when
- Creating or refactoring domain models
- Designing aggregates, entities, or value objects
- Implementing domain or application services
- Introducing new business concepts
- Working with invariants or transaction boundaries

## Do not use this skill when
- The task is purely UI or frontend
- The task is infrastructure-only (Docker, CI, etc.)
- The task is simple data transformation with no domain logic

---

## Required context
Before writing code, inspect:
- root `AGENTS.md`
- existing domain models in the bounded context
- repository interfaces and service patterns

---

## Core modeling philosophy

### Aggregate root
Use an aggregate root only when:
- it enforces invariants that must be consistent in a single transaction
- it is the entry point for modifying related entities
- external code should not modify internal entities directly

Do NOT:
- create aggregates just because something has an ID
- create large aggregates without strong invariants

---

### Entity
Use an entity when:
- identity matters
- it changes over time
- it belongs inside an aggregate

Do NOT:
- expose entity mutation outside the aggregate root

---

### Value object
Use a value object when:
- identity does not matter
- equality is based on value
- it can be immutable

Rules:
- must be immutable
- must validate itself at creation
- no repository
- no independent lifecycle

---

### Domain service
Use a domain service when:
- logic does not belong to a single aggregate
- it represents domain policy

Do NOT:
- move logic here if it can live on an aggregate
- use it as a dumping ground

---

### Application service
Use an application service to:
- orchestrate use cases
- load/save aggregates
- coordinate transactions

Do NOT:
- put core business logic here

---

## Required decision process (MANDATORY BEFORE CODING)

You MUST explicitly determine:

1. What business concept is being modeled?
2. What invariant must be protected?
3. What is the transaction boundary?
4. Does this concept require identity?
5. Can it be immutable?
6. Why is this NOT:
    - a value object?
    - an entity?
    - a domain service?
    - an aggregate root?

---

## Workflow

1. Identify the domain concept
2. Classify it:
    - aggregate root
    - entity
    - value object
    - domain service
    - application service
3. Define invariants
4. Define transaction boundary
5. Justify classification
6. THEN generate code

---

## Constraints

- Prefer rich domain models (behavior over data)
- Keep persistence concerns OUT of domain models
- Repositories do NOT own transactions
- Follow existing repo conventions strictly
- Avoid unnecessary abstractions
- Avoid anemic models

---

## Common mistakes to avoid

- Turning every concept into an aggregate root
- Mutable value objects
- Domain logic inside application services
- Creating domain services unnecessarily
- Splitting aggregates without considering invariants
- Letting entities be modified outside the aggregate

---

## Output format (REQUIRED)

### Modeling Decision
- Classification:
- Invariants:
- Transaction boundary:
- Why this design:

### Code
(implementation)

### Notes
- Tradeoffs
- Future considerations
