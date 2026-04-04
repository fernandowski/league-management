# Concurrency Invariants

This codebase uses database-enforced invariants for create commands, pessimistic locking for `League` membership mutations, and optimistic locking for `Season` state transitions.

## Aggregate Roots

`User`

- Invariant: one account per email.
- Locking choice: `UNIQUE` constraint on `users.email`.
- Why: registration is a create command keyed by a natural identifier; database uniqueness is the source of truth.

`Organization`

- Invariant: an owner cannot create two organizations with the same name.
- Locking choice: unique index on `(user_id, lower(name))`.
- Why: duplicate prevention belongs in the schema; application-level existence checks are advisory only.

`League`

- Invariant: an organization cannot create two leagues with the same name.
- Locking choice: unique index on `(organization_id, lower(name))`.
- Why: duplicate league creation is a create concern and should be enforced atomically by Postgres.

`League` membership set

- Invariant: membership add/remove commands must not lose unrelated membership changes, and the same team can appear only once in a league.
- Locking choice: pessimistic locking with `SELECT ... FOR UPDATE` on the league row, plus the existing unique index on `(league_id, team_id)`.
- Why: membership commands rewrite the aggregate child-entity set. Row locking serializes conflicting writers and avoids lost updates.

`Team`

- Invariant: team creation must atomically create the team row, organization link, and user-role rows.
- Locking choice: transaction boundary around the full repository save.
- Why: the consistency boundary spans multiple tables, so atomic commit matters more than row-level contention control.

`Season`

- Invariant: only one active season (`pending`, `planned`, `in_progress`, `paused`) may exist per league.
- Locking choice: partial unique index on `seasons(league_id)` for active statuses.
- Why: this is a cross-aggregate invariant and should be enforced centrally by the database.

`Season` state and match updates

- Invariant: season transitions and score updates must not silently overwrite each other.
- Locking choice: optimistic locking with `seasons.version`.
- Why: `Season` is the aggregate root for rounds and matches in this model. Version checks surface concurrent writes without serializing all reads.
