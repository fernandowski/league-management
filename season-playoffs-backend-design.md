# Season Playoffs Backend Design

This document translates the playoff planning into a backend design proposal aligned with the current `apps/api` architecture.

It assumes:

- `Season` remains the aggregate root
- playoffs are a phase inside the same season
- optimistic locking on `seasons.version` remains the primary aggregate concurrency mechanism
- the backend continues to follow the current DDD split:
  - domain owns business transitions
  - application services orchestrate
  - repositories rehydrate and persist

## Current Backend Baseline

Today the backend models:

- `Season` with `Status`, `Version`, `Rounds`
- regular-season round-robin schedule generation
- season start
- score updates
- round completion
- standings query

Current limitations relevant to playoffs:

- no explicit season phase beyond one status field
- no playoff rules model
- no standings snapshot
- no bracket concept
- no tie-level concept for multi-leg knockout logic
- no champion field

Current files that will be affected conceptually:

- [apps/api/internal/organization_management/domain/season.go](/home/fernando/ideaProjects/league-management/apps/api/internal/organization_management/domain/season.go:1)
- [apps/api/internal/organization_management/domain/matches.go](/home/fernando/ideaProjects/league-management/apps/api/internal/organization_management/domain/matches.go:1)
- [apps/api/internal/organization_management/application/services/season_service.go](/home/fernando/ideaProjects/league-management/apps/api/internal/organization_management/application/services/season_service.go:1)
- [apps/api/internal/organization_management/interfaces/http/controllers/seasons_controller.go](/home/fernando/ideaProjects/league-management/apps/api/internal/organization_management/interfaces/http/controllers/seasons_controller.go:1)
- [apps/api/internal/organization_management/infrastructure/repositories/season_repository.go](/home/fernando/ideaProjects/league-management/apps/api/internal/organization_management/infrastructure/repositories/season_repository.go:1)

## Design Goals

- keep playoff progression inside the `Season` aggregate
- model `PlayoffTie` explicitly so multi-leg aggregate scoring is correct
- keep rules typed and validated, not dynamic or script-based
- preserve invariants with the minimal concurrency mechanism
- avoid breaking current regular-season flows while extending them

## Aggregate Boundary

`Season` should remain the aggregate root.

Reason:

- playoff qualification depends on regular-season standings
- playoff seeding depends on season state
- playoff bracket generation is a season transition
- playoff tie resolution affects season progression
- champion belongs to the season

This means:

- no separate top-level `Playoff` aggregate for MVP
- no separate "postseason season"
- no repository methods that update playoff rows independently of season version checks for state-changing commands

## Proposed Domain Language

### Existing

- `Season`
- `Round`
- `Match`

### New

- `SeasonPhase`
- `PlayoffRules`
- `StandingsSnapshot`
- `QualifiedTeam`
- `PlayoffBracket`
- `PlayoffRound`
- `PlayoffTie`
- `PlayoffLeg`
- `AggregateScore`
- `TieResolutionRule`
- `Champion`

## Proposed Domain Shape

The current `Season` shape is too small for playoffs. Conceptually it should evolve toward:

```go
type Season struct {
    ID             string
    LeagueId       string
    Name           string
    Status         SeasonStatus
    Phase          SeasonPhase
    Version        int
    Rounds         []Round
    PlayoffRules   *PlayoffRules
    PlayoffBracket *PlayoffBracket
    ChampionTeamID *string
}
```

Suggested new types:

```go
type SeasonPhase string

const (
    SeasonPhaseRegularSeason SeasonPhase = "regular_season"
    SeasonPhasePlayoffs      SeasonPhase = "playoffs"
    SeasonPhaseCompleted     SeasonPhase = "completed"
)

type PlayoffRules struct {
    QualificationType            string
    QualifierCount               int
    ReseedEachRound              bool
    ThirdPlaceMatch              bool
    Rounds                       []PlayoffRoundRule
    AllowAdminSeedOverride       bool
}

type PlayoffRoundRule struct {
    Name                         string
    Legs                         int
    HigherSeedHostsSecondLeg     bool
    TiedAggregateResolution      TieResolutionRule
}

type PlayoffBracket struct {
    StandingsSnapshot StandingsSnapshot
    QualifiedTeams    []QualifiedTeam
    Rounds            []PlayoffRound
    WinnerTeamID      *string
}

type QualifiedTeam struct {
    TeamID string
    Seed   int
}

type PlayoffRound struct {
    Name  string
    Order int
    Ties  []PlayoffTie
}

type PlayoffTie struct {
    ID              string
    RoundName       string
    HomeSeed        int
    AwaySeed        int
    HomeTeamID      string
    AwayTeamID      string
    Legs            []PlayoffLeg
    WinnerTeamID    *string
    ResolutionUsed  *TieResolutionRule
    Status          PlayoffTieStatus
}

type PlayoffLeg struct {
    MatchID       string
    LegNumber     int
    HomeTeamID    string
    AwayTeamID    string
    HomeScore     int
    AwayScore     int
    Status        MatchStatus
}
```

This is a domain proposal, not necessarily the final persistence shape.

## Season Status Strategy

The current single `SeasonStatus` enum is not expressive enough once playoffs exist.

Recommended approach:

- keep `SeasonStatus` for lifecycle
- add `SeasonPhase` for competition stage

Suggested lifecycle meanings:

- `pending`
- `planned`
- `in_progress`
- `paused`
- `finished`

Suggested phase meanings:

- `regular_season`
- `playoffs`
- `completed`

Example:

- season in active league play:
  - `Status = in_progress`
  - `Phase = regular_season`

- season in playoff play:
  - `Status = in_progress`
  - `Phase = playoffs`

- season fully done:
  - `Status = finished`
  - `Phase = completed`

This avoids inventing a very wide status enum and keeps reasoning cleaner.

## Invariants And Protection

These are the core invariants and how to protect them.

### 1. One Playoff Ruleset Per Season

Invariant:

- a season has at most one active playoff rules configuration

Mechanism:

- database unique constraint on `season_playoff_rules.season_id`
- aggregate validation in `Season.ConfigurePlayoffRules`

### 2. Playoffs Cannot Start Before Regular Season Ends

Invariant:

- playoff bracket generation and playoff matches are invalid until regular season is complete

Mechanism:

- domain guard on `Season.GeneratePlayoffBracket`
- phase/status checks

### 3. Qualification Must Be Based On Frozen Standings

Invariant:

- playoff qualification and seeding must be derived from a fixed standings snapshot

Mechanism:

- create a `StandingsSnapshot` at regular-season completion
- persist one authoritative snapshot for playoff generation
- never compute qualifiers from a mutable live query after playoffs begin

### 4. Qualified Teams And Seeds Must Be Unique

Invariant:

- no duplicate team in playoff bracket
- no duplicate seed in playoff bracket

Mechanism:

- domain validation when generating bracket
- unique constraints on `(playoff_bracket_id, team_id)` and `(playoff_bracket_id, seed)` if persisted separately

### 5. Team Cannot Be In Two Ties In Same Round

Invariant:

- a team appears at most once per playoff round

Mechanism:

- bracket generation validation
- unique constraint per round/team if round-team rows are persisted

### 6. Tie Winner Comes From Tie, Not Match

Invariant:

- no playoff winner can be derived from one leg when rules require aggregate resolution

Mechanism:

- explicit `PlayoffTie` domain logic
- leg score update only mutates a leg and recomputes tie state
- winner can only be set through tie-resolution rules

### 7. Champion Can Only Come From Final Tie

Invariant:

- champion must be the resolved winner of the final playoff tie

Mechanism:

- domain guard in `Season.FinishSeason` or `Season.CompleteFinalTie`
- champion set only after final tie is resolved

### 8. Rules Cannot Change After Playoff Start

Invariant:

- playoff rules and seeds are immutable once playoffs begin

Mechanism:

- domain guard
- optimistic locking on season version

### 9. Concurrent Progression Must Not Overwrite State

Invariant:

- tie resolution, winner advancement, and champion assignment must not lose concurrent changes

Mechanism:

- optimistic locking on `seasons.version`

This matches the current repo rule for season aggregate state changes.

## Concurrency Decision

### Invariant Being Protected

- season progression must not silently overwrite playoff state

### Race

- two admins record scores or resolve ties concurrently
- one admin generates the bracket while another updates rules
- one admin finishes the season while another updates a playoff leg

### Minimal Mechanism

- optimistic locking using `seasons.version`

### Why This Is Better Than Alternatives

- playoff progression is aggregate state, not an independent child-collection rewrite
- pessimistic locking would be broader than needed
- the repo already uses season version checks, so this extends the current pattern cleanly

## Proposed Persistence Model

The persistence layer should distinguish:

- regular-season match rows
- playoff configuration rows
- playoff state rows

Suggested tables:

### `season_playoff_rules`

- `season_id` unique not null
- `qualification_type`
- `qualifier_count`
- `reseed_each_round`
- `third_place_match`
- `allow_admin_seed_override`
- `created_at`
- `updated_at`

### `season_playoff_round_rules`

- `id`
- `season_id`
- `round_order`
- `round_name`
- `legs`
- `higher_seed_hosts_second_leg`
- `tied_aggregate_resolution`

### `season_standings_snapshots`

- `id`
- `season_id`
- `snapshot_type`
- `created_at`

### `season_standings_snapshot_entries`

- `snapshot_id`
- `team_id`
- `rank`
- `points`
- `goal_difference`
- `goals_for`
- `goals_against`
- any other standing fields needed for auditability

### `season_playoff_brackets`

- `id`
- `season_id` unique not null
- `standings_snapshot_id`
- `status`
- `winner_team_id`

### `season_playoff_qualified_teams`

- `playoff_bracket_id`
- `team_id`
- `seed`

### `season_playoff_ties`

- `id`
- `playoff_bracket_id`
- `round_name`
- `round_order`
- `slot_order`
- `home_seed`
- `away_seed`
- `home_team_id`
- `away_team_id`
- `winner_team_id`
- `status`
- `resolution_used`
- `next_tie_id` nullable

### `season_playoff_legs`

- `id`
- `playoff_tie_id`
- `match_id`
- `leg_number`
- `home_team_id`
- `away_team_id`
- `home_score`
- `away_score`
- `status`

## Why Separate Playoff Tables Are Better Than Reusing `season_schedules`

Current `season_schedules` works for regular-season rounds, but it is too flat for playoffs.

Problems if playoffs are shoved into `season_schedules` only:

- no explicit tie boundary for aggregate scoring
- no clear winner advancement path
- awkward final/champion derivation
- awkward support for two-leg semantics
- hard-to-enforce invariants around one team per round/tie

Recommendation:

- keep `season_schedules` for regular season
- create explicit playoff tables
- optionally still use `match_id` as the underlying leg/match identifier

## Repository Changes

The current `SeasonRepository` loads season plus regular-season schedules in one query. That will need to expand.

Recommended repository responsibilities:

- `FindByID` rehydrates:
  - season row
  - regular-season rounds
  - playoff rules
  - playoff bracket
  - playoff ties
  - playoff legs

- `Save` persists:
  - season row with version check
  - playoff rules upsert
  - bracket rows upsert
  - ties and legs upsert

Important rule:

- state-changing playoff persistence should remain inside the same transaction as the season version update

Potential repository additions:

- `FetchPlayoffBracket(seasonID string)`
- `FetchPlayoffQualificationPreview(seasonID string)`
- `FetchSeasonSummary(seasonID string)`

Those are read models and do not need to load the full aggregate.

## Domain Methods To Add

These methods should live on `Season` or season-owned child types.

### Configuration

- `ConfigurePlayoffRules(rules PlayoffRules) error`

Guards:

- only before playoffs start
- only while season is editable
- rules must validate

### End Of Regular Season

- `CompleteRegularSeason(snapshot StandingsSnapshot) error`

Guards:

- all regular-season rules satisfied
- standings snapshot must be complete

Effects:

- freeze snapshot for qualification
- move season phase toward playoff setup

### Bracket Generation

- `GeneratePlayoffBracket() error`

Guards:

- playoff rules exist
- standings snapshot exists
- bracket not already generated

Effects:

- select qualified teams
- assign seeds
- build round/tie structure

### Playoff Start

- `StartPlayoffs() error`

Guards:

- bracket generated
- first-round ties ready

Effects:

- move phase to `playoffs`
- activate first playable legs

### Playoff Result Recording

- `RecordPlayoffLegScore(tieID string, legNumber int, homeScore, awayScore int) error`

Guards:

- season in playoff phase
- target tie and leg exist
- tie unresolved
- season not finished

Effects:

- updates leg score
- recomputes aggregate state
- may mark tie ready for resolution

### Tie Resolution

- `ResolvePlayoffTie(tieID string, input TieResolutionInput) error`

Guards:

- all required legs complete or rule allows manual resolution
- tie not already resolved

Effects:

- sets winner
- records resolution used
- advances winner if next tie can now be completed structurally

### Season Completion

- `FinishSeason() error`

Guards:

- final tie resolved
- champion determinable

Effects:

- champion set
- phase = completed
- status = finished

## Application Service Changes

`SeasonService` should gain explicit orchestration commands rather than generic update methods.

Recommended additions:

- `ConfigurePlayoffRules(orgOwnerID, seasonID string, dto ConfigurePlayoffRulesDTO) error`
- `GeneratePlayoffBracket(orgOwnerID, seasonID string) error`
- `StartPlayoffs(orgOwnerID, seasonID string) error`
- `RecordPlayoffLegScore(orgOwnerID, seasonID string, dto RecordPlayoffLegScoreDTO) error`
- `ResolvePlayoffTie(orgOwnerID, seasonID string, dto ResolvePlayoffTieDTO) error`
- `PlayoffBracket(orgOwnerID, seasonID string) (map[string]interface{}, error)`
- `PlayoffQualificationPreview(orgOwnerID, seasonID string) (map[string]interface{}, error)`
- `SeasonSummary(orgOwnerID, seasonID string) (map[string]interface{}, error)`

Authorization should stay where it is now:

- service loads season
- service loads league
- service loads organization
- owner authorization is checked before command execution

## HTTP API Proposal

These are proposal-level routes, not final route commitments.

### Commands

- `PUT /seasons/{season_id}/playoffs/rules`
- `POST /seasons/{season_id}/playoffs/bracket/generate`
- `POST /seasons/{season_id}/playoffs/start`
- `PUT /seasons/{season_id}/playoffs/ties/{tie_id}/legs/{leg_number}/score`
- `POST /seasons/{season_id}/playoffs/ties/{tie_id}/resolve`

### Queries

- `GET /seasons/{season_id}/playoffs`
- `GET /seasons/{season_id}/playoffs/qualification-preview`
- `GET /seasons/{season_id}/summary`

## DTOs To Add

Suggested DTOs:

### `ConfigurePlayoffRulesDTO`

- `qualifier_count`
- `reseed_each_round`
- `third_place_match`
- `rounds`

### `PlayoffRoundRuleDTO`

- `name`
- `legs`
- `higher_seed_hosts_second_leg`
- `tied_aggregate_resolution`

### `RecordPlayoffLegScoreDTO`

- `home_score`
- `away_score`

### `ResolvePlayoffTieDTO`

- `resolution_type`
- `winner_team_id`
- `notes`

The exact DTO naming should stay consistent with the existing codebase style.

## Migration Strategy

To reduce rollout risk, do this in phases.

### Phase 1: Schema Preparation

- add `season_phase` to `seasons`
- add `champion_team_id` to `seasons`
- add playoff rules tables
- add standings snapshot tables
- add playoff bracket/tie/leg tables

### Phase 2: Domain Model

- extend `Season`
- add playoff value objects/entities
- add rule validation

### Phase 3: Repository Rehydration

- update `FindByID`
- update `Save`
- add read-model queries

### Phase 4: Commands And Routes

- configure rules
- generate bracket
- start playoffs
- record leg score
- resolve tie

### Phase 5: Tests

- domain tests
- repository tests
- concurrency tests
- API/service tests

## Test Plan

### Domain Tests

- valid playoff rules can be configured before season start
- playoff rules cannot change after playoff phase begins
- top `N` qualification produces unique teams and seeds
- two-leg tie computes correct aggregate score
- tie cannot resolve early
- final tie winner becomes champion

### Service Tests

- unauthorized user cannot configure rules
- unauthorized user cannot generate bracket
- season cannot start playoffs without bracket
- score updates after season finish are rejected

### Concurrency Tests

- concurrent bracket generation yields one success and one optimistic conflict
- concurrent tie resolution yields one success and one optimistic conflict
- concurrent finalization and score update yields one success and one conflict

### Repository Tests

- full aggregate rehydrates regular season plus playoffs correctly
- rules/qualified teams/ties/legs persist consistently
- uniqueness constraints reject invalid duplicates

## Recommended MVP Scope

For first implementation, keep the backend scope narrow:

- one qualification type: `top_n`
- `4` or `8` qualifiers only
- fixed bracket only
- one-leg final
- optional two-leg semifinal
- no away-goals rule
- penalties as the only tied aggregate resolution
- no reseeding after bracket generation
- no admin seed override after playoff generation

This avoids overbuilding while still supporting real playoff behavior.

## Main Risks

- overloading current regular-season models instead of introducing `PlayoffTie`
- treating standings as a live query instead of freezing a snapshot
- letting repositories mutate playoff state independently of season version checks
- allowing rules or seeds to change after playoff start

## Summary

The backend should evolve by keeping `Season` as the aggregate root and introducing explicit playoff concepts under it:

- `PlayoffRules`
- `StandingsSnapshot`
- `PlayoffBracket`
- `PlayoffTie`
- `PlayoffLeg`

The most important technical decisions are:

- use a separate `SeasonPhase`
- freeze standings before qualification
- model ties explicitly for aggregate scoring
- protect state changes with optimistic locking on `seasons.version`

That gives a design that fits the current architecture and preserves the invariants that matter once playoffs exist.
