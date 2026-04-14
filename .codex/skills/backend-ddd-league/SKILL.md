---
name: backend-ddd-league
description: Use when working in `apps/api` on Go backend features, bug fixes, refactors, or reviews that must follow this repo's domain-driven design, aggregate boundaries, concurrency rules, and testing expectations. Trigger for changes to domain models, application services, repositories, controllers, DTOs, routes, migrations, and backend tests.
---

# Backend DDD League

## Overview

Use this skill for backend work in `apps/api`. It encodes this repo's architectural boundaries, DDD expectations, and concurrency rules so changes stay inside the correct layer and invariants remain explicit.

Start by reading `apps/api/AGENTS.md` and the existing feature slice you are modifying before proposing code.

## Workflow

1. Identify the user-visible command or query.
2. Identify the aggregate and invariant being changed.
3. Decide the minimal mechanism required to protect that invariant.
4. Change the smallest set of layers needed.
5. Add or update tests that prove the invariant and behavior.

For any concurrency-sensitive change, explicitly answer:
1. What invariant is being protected?
2. What race can violate it?
3. What is the minimal mechanism?
4. Why is that mechanism better than alternatives here?

## Architecture Rules

### Controllers

- Parse request data, authenticate, call application services, map errors to HTTP.
- Do not contain business rules.
- Do not decide locking or concurrency strategy.

### Application Services

- Coordinate use cases and transaction boundaries.
- Load aggregates through repository interfaces.
- Enforce authorization and orchestration across aggregates.
- Do not embed SQL or persistence details.
- Do not duplicate domain invariants that belong in aggregates.

### Domain

- Keep domain types pure Go with no HTTP, SQL, ORM, or framework concerns.
- Put business invariants and state transitions on aggregates and entities.
- Prefer explicit methods such as `Start`, `ChangeMatchScore`, `CompleteCurrentRound` over mutating fields from services.
- If behavior spans multiple aggregates but is still domain logic, consider `domain/domainservices`.

### Repositories

- Persist and rehydrate aggregates faithfully.
- Own optimistic or pessimistic locking details.
- Translate database errors into domain or application-level conflicts where appropriate.
- If a new domain field affects behavior, update both load and save paths.

## Repo-Specific DDD Expectations

- `Season` is an aggregate root. Match updates and season state transitions must go through `Season` behavior, not ad hoc repository updates.
- Child collections that can lose updates, such as league memberships, use parent-row locking when mutated.
- Aggregate state conflicts use optimistic locking via `version` checks.
- Database constraints are authoritative for uniqueness and active-season invariants.
- Application checks may improve UX, but they do not replace schema constraints.

## Change Patterns

### New command

- Add or extend a domain method first if state changes are involved.
- Wire it through the relevant application service.
- Expose it via controller and route only after domain behavior is clear.
- Update repository persistence for any new fields or transitions.

### New query or read model

- Keep it out of the aggregate unless it is required for invariant enforcement.
- Add repository fetch/projection methods for read-optimized responses.
- Avoid polluting domain entities with view-only data.

### New invariant

- State it explicitly in code review notes and tests.
- Prefer database enforcement for uniqueness or cross-row constraints.
- Prefer pessimistic locking for child-collection rewrites.
- Prefer optimistic locking for aggregate state transitions that should fail fast.

## Concurrency Defaults

Use the rules from `apps/api/AGENTS.md`:

- Uniqueness invariants: database unique constraints/indexes.
- Cross-aggregate active season rule: database partial unique index.
- Child collection rewrites: transaction plus `SELECT ... FOR UPDATE` on the parent row.
- Aggregate state changes: optimistic locking on the aggregate version.
- Do not add broader locks unless the invariant requires them.
- Do not silently overwrite concurrent updates.

## File Targets

Typical files for backend changes:

- Domain: `apps/api/internal/*/domain/*.go`
- Domain services: `apps/api/internal/*/domain/domainservices/*.go`
- App services: `apps/api/internal/*/application/services/*.go`
- HTTP controllers/routes: `apps/api/internal/*/interfaces/http/**/*`
- Persistence: `apps/api/internal/*/infrastructure/repositories/*.go`
- DTOs: `apps/api/internal/shared/dtos/*.go`
- Migrations: `apps/api/migrations/*.sql`

## Testing Rules

When backend behavior changes, update the narrowest useful tests:

- Domain rules/state transitions: domain tests.
- Service orchestration/locking/conflicts: service tests.
- Concurrency behavior: dedicated concurrency tests.
- Persistence invariants or constraints: repository/integration tests when needed.

For concurrency-related changes, tests must prove:

- duplicate creation is blocked when required
- concurrent writes do not lose updates
- optimistic conflicts fail explicitly
- final persisted state matches the invariant

## Review Checklist

Before finishing backend work, verify:

- Domain logic lives in domain methods, not controllers.
- Services coordinate; repositories persist.
- New fields are rehydrated and saved consistently.
- Invariants are backed by the correct mechanism.
- Tests cover the new behavior and any concurrency risk.
- Unrelated files were not changed without reason.
