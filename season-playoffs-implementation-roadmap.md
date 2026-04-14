# Season Playoffs Implementation Roadmap

This roadmap combines the backend and frontend playoff planning into one execution plan with logical tasks, dependencies, and implementation order.

It is designed so the work can be implemented incrementally in this repo without breaking the current season flow.

Related planning documents:

- [season-playoffs.md](/home/fernando/ideaProjects/league-management/season-playoffs.md:1)
- [season-playoffs-backend-design.md](/home/fernando/ideaProjects/league-management/season-playoffs-backend-design.md:1)
- [season-playoffs-frontend-plan.md](/home/fernando/ideaProjects/league-management/season-playoffs-frontend-plan.md:1)

## Product Scope For Implementation

Target MVP:

- playoffs are a phase inside the same season
- playoff rules are configurable before playoffs start
- qualification is based on frozen regular-season standings
- bracket generation is deterministic
- one-leg and two-leg ties are supported where configured
- aggregate scoring is computed at tie level
- a champion is produced at season completion
- frontend supports configuration, qualification preview, bracket view, playoff score entry, and champion display

Out of scope for first pass:

- public playoff pages
- manual seed drag-and-drop
- custom scriptable rules
- away goals
- rich animated bracket visualizations

## Guiding Constraints

- `Season` remains the aggregate root
- playoff state changes are protected by optimistic locking on `seasons.version`
- database constraints protect uniqueness and structural invariants
- backend state-changing commands land before frontend interactive flows depending on them
- frontend reuses the current season surface and adds a `Playoffs` tab

## Delivery Strategy

Build this in vertical slices, but in the right dependency order:

1. backend schema and domain foundation
2. backend read models and commands
3. frontend type and screen scaffolding
4. frontend interactive playoff flows
5. hardening with tests and UX refinement

## Milestone 1: Season Foundation

Goal:

- make the existing season model aware of playoff lifecycle without adding bracket behavior yet

### Backend Tasks

#### Task 1.1: Add season phase support

Outcome:

- `Season` can distinguish regular season, playoffs, and completed state

Scope:

- add `SeasonPhase` in domain
- add `season_phase` column in `seasons`
- update repository load/save
- preserve current behavior for existing seasons

Files likely touched:

- `apps/api/internal/organization_management/domain/season.go`
- `apps/api/internal/organization_management/infrastructure/repositories/season_repository.go`
- new migration in `apps/api/migrations`

Dependency:

- none

#### Task 1.2: Add champion support to season

Outcome:

- season can persist and expose final champion

Scope:

- add `champion_team_id` to `seasons`
- update aggregate, repository, and season detail read model

Dependency:

- Task 1.1

### Frontend Tasks

#### Task 1.3: Extend client season types

Outcome:

- mobile app can consume season phase and champion information safely

Scope:

- update `useData.ts` response types
- update season details typing

Files likely touched:

- `apps/mobile/hooks/useData.ts`

Dependency:

- backend Tasks 1.1 and 1.2 for final API shape

#### Task 1.4: Update season overview to show phase/champion placeholders

Outcome:

- season UI can display new state as backend fields arrive

Scope:

- extend `SeasonInformation.tsx`
- keep layout backwards compatible

Dependency:

- Task 1.3

## Milestone 2: Playoff Rules Foundation

Goal:

- configure playoff rules before playoffs start

### Backend Tasks

#### Task 2.1: Add playoff rules persistence

Outcome:

- one ruleset can be stored per season

Scope:

- create `season_playoff_rules`
- create `season_playoff_round_rules`
- add unique constraint on `season_id`

Dependency:

- Task 1.1

#### Task 2.2: Add playoff rules domain model and validation

Outcome:

- backend can validate playoff configuration in domain language

Scope:

- add `PlayoffRules`
- add `PlayoffRoundRule`
- add `ConfigurePlayoffRules` method on `Season`
- enforce:
  - qualifier count valid
  - round definitions valid
  - rules immutable after playoffs start

Dependency:

- Task 2.1

#### Task 2.3: Add configure playoff rules service and route

Outcome:

- org owner can configure playoff rules through API

Scope:

- DTO for playoff rules
- service orchestration
- controller action
- route registration

Dependency:

- Task 2.2

#### Task 2.4: Add playoff rules read model

Outcome:

- frontend can fetch configured rules and summary state

Scope:

- season details summary includes playoff summary, or add dedicated playoff rules endpoint

Dependency:

- Task 2.3

### Frontend Tasks

#### Task 2.5: Add playoff types and rules DTO support

Outcome:

- mobile app can fetch and submit playoff rule configuration

Scope:

- add types for playoff rules and summaries
- add request payload shape

Dependency:

- Task 2.4

#### Task 2.6: Add `Playoffs` tab shell

Outcome:

- season UI has a dedicated place for playoff work

Scope:

- update `SeasonManagement.tsx`
- add a new `SeasonPlayoffs` container component
- show empty states for missing rules

Dependency:

- Task 2.5

#### Task 2.7: Build playoff rules form

Outcome:

- admin can configure playoff rules from the app

Scope:

- create `PlayoffRulesForm`
- save rules
- show read-only state after playoffs start

Dependency:

- Task 2.6

## Milestone 3: Standings Snapshot and Qualification

Goal:

- freeze regular-season standings and determine qualifiers safely

### Backend Tasks

#### Task 3.1: Add standings snapshot persistence

Outcome:

- playoff qualification can be based on frozen standings

Scope:

- create `season_standings_snapshots`
- create `season_standings_snapshot_entries`

Dependency:

- Task 1.1

#### Task 3.2: Add regular-season completion flow

Outcome:

- backend can formally close regular season and store standings snapshot

Scope:

- add `CompleteRegularSeason(snapshot)` semantics
- decide whether this extends current `CompleteCurrentRound` behavior or adds a new explicit completion step

Dependency:

- Task 3.1

#### Task 3.3: Add qualification preview logic

Outcome:

- backend can compute top-`N` qualifiers and seeds from snapshot

Scope:

- add qualification projection/read model
- enforce unique teams and seeds

Dependency:

- Task 2.2
- Task 3.2

#### Task 3.4: Add qualification preview endpoint

Outcome:

- frontend can show who qualifies before bracket generation

Dependency:

- Task 3.3

### Frontend Tasks

#### Task 3.5: Update standings screen for playoff context

Outcome:

- standings screen can show rank, qualified markers, and frozen state

Scope:

- extend `SeasonStanding.tsx`
- likely extend `TableList` column definitions only, no new primitive needed

Dependency:

- Task 3.4 or season summary support for qualified markers

#### Task 3.6: Build qualification preview section

Outcome:

- playoffs tab shows qualifiers and seeds

Scope:

- create `PlayoffQualificationPreview`
- show cut line and generate-bracket CTA

Dependency:

- Task 3.4

## Milestone 4: Bracket Generation

Goal:

- generate playoff bracket from rules plus standings snapshot

### Backend Tasks

#### Task 4.1: Add playoff bracket persistence

Outcome:

- one bracket can be stored per season

Scope:

- create `season_playoff_brackets`
- create `season_playoff_qualified_teams`
- add unique constraints for:
  - one bracket per season
  - unique team per bracket
  - unique seed per bracket

Dependency:

- Task 2.1
- Task 3.1

#### Task 4.2: Add playoff tie and leg persistence

Outcome:

- backend can persist explicit knockout ties and legs

Scope:

- create `season_playoff_ties`
- create `season_playoff_legs`
- add structural uniqueness constraints

Dependency:

- Task 4.1

#### Task 4.3: Add bracket generation domain logic

Outcome:

- season can generate deterministic bracket from rules and snapshot

Scope:

- add `PlayoffBracket`
- add `PlayoffRound`
- add `PlayoffTie`
- add `PlayoffLeg`
- add `GeneratePlayoffBracket`

Dependency:

- Task 2.2
- Task 3.3
- Task 4.2

#### Task 4.4: Add bracket generation service and route

Outcome:

- org owner can generate bracket through API

Dependency:

- Task 4.3

#### Task 4.5: Add playoff bracket read model

Outcome:

- frontend can render bracket without loading internal aggregate shape directly

Dependency:

- Task 4.3

### Frontend Tasks

#### Task 4.6: Add bracket response types

Outcome:

- mobile app can consume bracket payloads safely

Dependency:

- Task 4.5

#### Task 4.7: Build bracket view

Outcome:

- playoffs tab can render rounds, ties, and legs

Scope:

- create:
  - `PlayoffBracketView`
  - `PlayoffRoundSection`
  - `PlayoffTieCard`
  - `PlayoffLegRow`

Dependency:

- Task 4.6

#### Task 4.8: Add bracket generation action to UI

Outcome:

- admin can generate bracket from qualification preview

Dependency:

- Task 4.4
- Task 4.7

## Milestone 5: Playoff Match Management

Goal:

- allow score entry for playoff legs and correct aggregate computation

### Backend Tasks

#### Task 5.1: Add playoff start transition

Outcome:

- season can move from regular-season-complete to playoffs-in-progress

Scope:

- add `StartPlayoffs` if this is a distinct step
- otherwise define bracket generation as implicitly readying playoffs

Dependency:

- Task 4.3

#### Task 5.2: Add playoff leg score update command

Outcome:

- backend can record a score for one playoff leg

Scope:

- add `RecordPlayoffLegScore`
- recompute tie aggregate after save
- enforce that tie winner is not set early

Dependency:

- Task 4.3

#### Task 5.3: Add tie resolution logic

Outcome:

- backend can resolve one-leg or two-leg ties correctly

Scope:

- penalties as MVP tie-break resolution
- no away-goals support

Dependency:

- Task 5.2

#### Task 5.4: Add advancement logic

Outcome:

- winner of one tie populates the next round

Dependency:

- Task 5.3

#### Task 5.5: Add playoff score update route and bracket refresh query

Outcome:

- frontend can edit scores and see resulting tie/bracket state

Dependency:

- Task 5.2
- Task 5.5 typo intentionally avoided by using Task 5.4 as logical prerequisite

### Frontend Tasks

#### Task 5.6: Reuse score entry modal for playoff legs

Outcome:

- admin can edit playoff leg scores without inventing a new interaction pattern

Scope:

- extend `ManageMatchScoreModal.tsx` or wrap it for playoff use

Dependency:

- Task 5.5

#### Task 5.7: Add editable playoff leg actions in bracket view

Outcome:

- playable legs can be updated from the bracket screen

Dependency:

- Task 5.6

#### Task 5.8: Show aggregate and winner state in UI

Outcome:

- users can see tie winner, unresolved tie, and aggregate score clearly

Dependency:

- Task 5.7

## Milestone 6: Champion and Finished Season

Goal:

- finish season cleanly and show the winner

### Backend Tasks

#### Task 6.1: Add final tie completion and champion assignment

Outcome:

- season gets a champion from the final playoff tie

Dependency:

- Task 5.4

#### Task 6.2: Add finish season transition for playoff flow

Outcome:

- season ends with:
  - `status = finished`
  - `phase = completed`
  - `champion_team_id` set

Dependency:

- Task 6.1

#### Task 6.3: Add season summary read model

Outcome:

- frontend can render final season state, champion, and playoff summary cleanly

Dependency:

- Task 6.2

### Frontend Tasks

#### Task 6.4: Add champion card and finished-state summary

Outcome:

- season overview and playoffs tab both reflect completion

Dependency:

- Task 6.3

#### Task 6.5: Lock completed UI actions

Outcome:

- score entry and rules editing are visibly disabled after completion

Dependency:

- Task 6.4

## Milestone 7: Hardening

Goal:

- prove invariants and reduce regression risk

### Backend Tasks

#### Task 7.1: Domain tests for playoff rules and tie resolution

#### Task 7.2: Repository tests for full aggregate rehydration

#### Task 7.3: Concurrency tests

Scenarios:

- concurrent rule update vs playoff start
- concurrent bracket generation
- concurrent leg score update and tie resolution
- concurrent finalization and score update

### Frontend Tasks

#### Task 7.4: Empty, loading, and error states for every playoff phase

#### Task 7.5: UX polish for mobile layout and large-screen layout

#### Task 7.6: Narrow verification on affected season screens

## Recommended First Working Slice

If we want the fastest meaningful implementation path, start here:

1. Task 1.1: season phase
2. Task 2.1: playoff rules tables
3. Task 2.2: playoff rules domain model
4. Task 2.3: configure playoff rules API
5. Task 2.6: playoffs tab shell
6. Task 2.7: playoff rules form

That slice gives the product a visible playoff setup workflow without yet needing bracket generation.

## Recommended Second Slice

1. Task 3.1: standings snapshot
2. Task 3.3: qualification preview logic
3. Task 3.4: qualification preview endpoint
4. Task 3.6: qualification preview UI

That makes the season-to-playoffs handoff real.

## Recommended Third Slice

1. Task 4.1: bracket persistence
2. Task 4.2: tie and leg persistence
3. Task 4.3: bracket generation logic
4. Task 4.4: bracket generation API
5. Task 4.7: bracket view
6. Task 4.8: generate bracket action

That completes the bracket foundation.

## Recommended Fourth Slice

1. Task 5.2: playoff leg score update
2. Task 5.3: tie resolution
3. Task 5.4: advancement
4. Task 5.6: playoff score entry UI
5. Task 5.8: aggregate and winner state
6. Task 6.1: champion assignment
7. Task 6.4: champion display

That completes the competition loop.

## Implementation Notes For Working Sessions

To keep delivery clean, each implementation session should try to stay inside one of these categories:

- schema + repository plumbing
- domain behavior + tests
- service/controller wiring
- frontend read-only UI
- frontend interactive UI

Do not mix broad backend schema changes with broad frontend state work in the same pass unless the backend contract is already stable.

## Task Ownership Recommendation

If implemented sequentially by one agent, use this order:

- backend foundation first
- frontend shell second
- backend bracket logic third
- frontend bracket rendering fourth
- backend result progression fifth
- frontend interaction sixth
- hardening last

If split among multiple workers later, the safest concurrent split is:

- worker A: schema + repository
- worker B: domain + service logic
- worker C: frontend shell + view rendering

Only after API contracts stabilize.

## Summary

The implementation should not begin with bracket visuals. The logical order is:

1. make season phase-aware
2. add playoff rules
3. freeze standings and preview qualification
4. generate bracket
5. record playoff legs and resolve ties
6. assign champion and finish season
7. polish and harden

This gives us a roadmap that can be executed incrementally while preserving aggregate boundaries and frontend clarity.
