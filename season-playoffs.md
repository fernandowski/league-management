# Season Playoffs Planning

This document defines the planned MVP feature for playoffs at the season level.

The current backend supports a regular season with round-robin scheduling, standings, and linear round progression. Playoffs should extend that model as a post-regular-season phase of the same season, with qualification, seeding, knockout rounds, and multi-leg ties.

## Goals

- Allow a season to declare playoff rules before the season starts
- Determine which teams qualify when the regular season ends
- Seed qualified teams from final standings
- Support knockout rounds with one-leg or two-leg ties
- Support home-and-away aggregate scoring across multiple matches
- Produce a champion at the end of the season
- Keep the rules explicit and auditable

## Non-Goals For MVP

- Full custom competition scripting
- Group-stage playoffs
- Dynamic bracket editing after playoff start
- Support for every regional tiebreak rule from day one
- Live bracket visualization concerns in this document

## Product Questions To Settle

These need explicit product decisions because they change backend design:

- How many teams qualify: `2`, `4`, `6`, `8`, or configurable?
- Is qualification always based on regular-season standings?
- Are playoff seeds fixed from standings, or can admins override them before playoff start?
- Are rounds single-leg, two-leg, or configurable per round?
- Does the final also support two legs, or is it always single-leg?
- Is there a third-place match?
- Are away goals used?
- If aggregate is tied, what resolves it:
  extra time, penalties, replay, higher seed advances, or manual admin decision?
- Can an admin manually adjust a playoff result after a match is finalized?

## Recommended MVP Decisions

To keep the first version implementable, use these defaults:

- Qualification comes from final regular-season standings
- Top `N` teams qualify, where `N` is configurable per season
- Seed order follows final standings
- Bracket is fixed once playoffs start
- Semifinals can be one-leg or two-leg
- Final is single-leg by default
- No away-goals rule in MVP
- If a tie is level on aggregate after all legs, resolve by penalties
- Third-place match is optional but not required for MVP
- Admins may regenerate the playoff bracket only before the playoff phase starts

This gives enough flexibility without requiring a fully open-ended rules engine on the first pass.

## Why This Should Be Modeled As Rules

Playoffs are not a single hardcoded format. Different leagues vary in:

- number of qualifiers
- seeding behavior
- bracket shape
- whether higher seeds get byes
- whether a round is one-leg or two-leg
- how tied aggregate scores are resolved

If these decisions are scattered across services and queries, the model will become brittle quickly. The backend should treat playoff configuration as season-owned rules, and playoff generation should be derived from those rules.

For MVP, this does not need a generic scripting engine. It can be a typed rules model with explicit fields and validation. The phrase "rules engine" here should mean "configurable, validated rule objects and deterministic evaluators", not arbitrary user-authored code.

## Ubiquitous Language

These terms should be used consistently in backend code, API payloads, and planning discussions.

- `Season`: the full competition container for one league competition cycle
- `SeasonPhase`: the current stage of the season
- `Regular Season`: the standings-driven phase before playoffs
- `Playoffs`: the knockout phase inside the same season
- `Playoff Rules`: season-owned configuration that defines qualification, seeding, bracket shape, and tie resolution
- `Qualified Team`: a team that earned entry into the playoffs from regular-season results
- `Seed`: a team’s playoff rank derived from final regular-season standings
- `Bracket`: the generated playoff structure for the season
- `Playoff Round`: a stage in the bracket such as semifinal or final
- `Playoff Tie`: a contest between two teams that may contain one or more legs
- `Leg`: an individual match inside a playoff tie
- `Aggregate Score`: the combined score across all legs in a playoff tie
- `Tie Resolution`: the rule used when aggregate score is level
- `Champion`: the winner of the final playoff tie, or the regular-season winner if playoffs are disabled
- `Standings Snapshot`: the frozen ranking used to determine playoff qualification and seeding

Recommended language to avoid:

- avoid calling playoffs a "new season"
- avoid calling a two-leg tie a "round" when the actual concept is a `PlayoffTie`
- avoid using `match winner` when the meaningful outcome is `tie winner`

## Core Framing

Playoffs should not be modeled as a new season.

The correct domain framing is:

- one season
- multiple phases inside that season

That avoids splitting one competition year into two aggregates and keeps these concerns together:

- regular-season standings
- qualification
- seeding
- playoff bracket
- champion

Suggested season phases:

- regular season
- playoffs
- finished

The playoff phase should be configured before the season starts and activated only after the regular season is completed.

## Domain Concept

Suggested concepts:

- `Season`
- `SeasonPhase`
- `StandingsSnapshot`
- `PlayoffRules`
- `PlayoffBracket`
- `PlayoffRound`
- `PlayoffTie`
- `PlayoffLeg`
- `QualificationRule`
- `TieBreakerRule`

## Core Rules To Model

### 1. Qualification Rules

Defines who qualifies for playoffs.

Suggested fields:

- `qualification_type`: `top_n`
- `qualifier_count`
- `use_final_regular_season_standings`: boolean
- `allow_admin_override_before_start`: boolean

MVP recommendation:

- only support `top_n`

Future options:

- top teams by conference/group
- host automatic qualification
- wildcard qualifiers

### 2. Seeding Rules

Defines how qualified teams are ordered into the bracket.

Suggested fields:

- `seed_source`: `standings`
- `reseed_each_round`: boolean
- `protect_higher_seed`: boolean

MVP recommendation:

- seed directly from final standings
- no reseeding after bracket creation

### 3. Bracket Rules

Defines the shape of the knockout stage.

Suggested fields:

- `bracket_type`: `fixed_knockout`
- `rounds`: array of round definitions
- `third_place_match`: boolean

Each round definition may include:

- `name`: `quarterfinal`, `semifinal`, `final`
- `tie_count`
- `legs`
- `higher_seed_hosts_second_leg`: boolean
- `single_match_host_rule`

### 4. Aggregate Tie Rules

Defines how multi-leg ties are decided.

Suggested fields:

- `aggregate_scoring_enabled`: boolean
- `away_goals_enabled`: boolean
- `tied_aggregate_resolution`

Suggested enum for `tied_aggregate_resolution`:

- `penalties`
- `extra_time_then_penalties`
- `higher_seed_advances`
- `admin_decision`

MVP recommendation:

- aggregate scoring on for two-leg ties
- away goals off
- penalties for tied aggregate

## Aggregate Boundary

The current codebase already treats `Season` as the aggregate root. That should remain true.

Recommended ownership:

- `Season` owns:
  - phase
  - regular-season rounds and matches
  - playoff rules
  - standings snapshot reference or embedded snapshot metadata
  - playoff bracket state
  - champion

- `PlayoffTie` and `PlayoffLeg` are children inside the `Season` aggregate, not separate top-level aggregates for MVP

Why:

- playoff qualification depends on season standings
- playoff generation depends on season rules
- tie advancement changes season state
- champion belongs to the season

This keeps all critical season progression invariants inside one aggregate boundary.

## Proposed Season Phases

The existing status model likely needs to evolve from a simple lifecycle state into an explicit phase-aware lifecycle.

Suggested distinction:

- `SeasonLifecycleStatus`
  - `pending`
  - `planned`
  - `in_progress`
  - `paused`
  - `finished`

- `SeasonPhase`
  - `regular_season`
  - `playoffs`
  - `completed`

Another acceptable option is one status enum that includes playoff-aware values, for example:

- `pending`
- `planned`
- `regular_season_in_progress`
- `regular_season_paused`
- `playoffs_pending`
- `playoffs_in_progress`
- `playoffs_paused`
- `finished`

The first option is cleaner conceptually because lifecycle and competition phase are different concerns.

## Proposed Aggregate Shape

For MVP, a season can conceptually look like this:

- `Season`
  - `ID`
  - `LeagueID`
  - `Name`
  - `LifecycleStatus`
  - `Phase`
  - `Version`
  - `RegularSeason`
  - `PlayoffRules`
  - `PlayoffBracket`
  - `ChampionTeamID`

- `RegularSeason`
  - `Rounds`

- `PlayoffBracket`
  - `StandingsSnapshotTaken`
  - `QualifiedTeams`
  - `Rounds`
  - `WinnerTeamID`

- `PlayoffRound`
  - `Name`
  - `Order`
  - `Ties`

- `PlayoffTie`
  - `ID`
  - `RoundName`
  - `HomeSeed`
  - `AwaySeed`
  - `HomeTeamID`
  - `AwayTeamID`
  - `Legs`
  - `WinnerTeamID`
  - `Status`
  - `Resolution`

- `PlayoffLeg`
  - `MatchID`
  - `LegNumber`
  - `HomeTeamID`
  - `AwayTeamID`
  - `Score`
  - `Status`

The backend does not need to persist this exact shape as JSON or as table names, but these are the domain concepts that should exist explicitly.

## Multi-Leg Head-to-Head Rules

This is the most important new capability.

A playoff tie may contain two matches:

- Leg 1: Team A vs Team B
- Leg 2: Team B vs Team A

The winner is determined from the combined score across both legs.

Example:

- Leg 1: Team A 2-1 Team B
- Leg 2: Team B 1-0 Team A
- Aggregate: Team A 2-2 Team B

The system then applies the configured aggregate tiebreak rule:

- if away goals enabled, compare away goals
- otherwise move to penalties or another configured rule

Important backend requirement:

The winner must not be derived from a single match row. It must be derived from a `PlayoffTie` aggregate over all legs in that tie.

## Core Commands

These are the meaningful state-changing operations the domain should expose.

- `ConfigurePlayoffRules`
- `StartRegularSeason`
- `CompleteRegularSeason`
- `GeneratePlayoffBracket`
- `StartPlayoffs`
- `SchedulePlayoffLeg`
- `RecordPlayoffLegScore`
- `ResolvePlayoffTie`
- `AdvancePlayoffWinners`
- `FinishSeason`

Avoid generic mutation-style commands such as:

- `UpdateSeason`
- `SaveBracket`
- `PatchPlayoff`

The commands should express business transitions, not persistence intent.

## Core Queries

Read models can stay separate from the aggregate, but the business language should be explicit.

- `GetSeasonSummary`
- `GetRegularSeasonStandings`
- `GetPlayoffQualificationPreview`
- `GetPlayoffBracket`
- `GetPlayoffTieDetails`
- `GetChampionSummary`

## Invariants

These are the key rules that should be protected explicitly in the backend.

### Season Ownership Invariants

- A season owns at most one playoff ruleset
- A season owns at most one playoff bracket
- A champion belongs to exactly one season outcome

### Phase Invariants

- Playoffs cannot begin before the regular season is complete
- Qualification cannot be computed before final regular-season standings are frozen
- Once playoffs start, playoff rules cannot change
- Once the season is finished, neither regular-season nor playoff scores may change

### Qualification Invariants

- Qualified teams must be derived from the frozen standings snapshot or an explicit allowed admin override before playoff start
- A team cannot qualify more than once
- The number of qualified teams must match playoff rules
- A non-qualified team cannot appear in the playoff bracket

### Seeding Invariants

- Each qualified team has one unique seed
- Seed ordering is deterministic from the standings snapshot unless an explicit pre-start override exists
- Seeding cannot change after playoff start

### Bracket Invariants

- The bracket shape must match the configured playoff rules
- Each playoff tie must contain exactly two teams unless a bye rule is explicitly supported
- A team cannot appear in two different active playoff ties in the same round
- Advancement paths must be deterministic and unique

### Multi-Leg Tie Invariants

- A playoff tie winner cannot be declared until all required legs are completed, unless the configured resolution rule permits an administrative decision
- Aggregate score must be computed from all completed legs in the tie
- A leg result changes tie state, but does not directly bypass tie-resolution rules
- The final tie winner is the only valid source of the season champion

### Champion Invariants

- A champion cannot be assigned before the final playoff tie is resolved
- Once champion is set, season status must move to finished
- Champion must be one of the qualified playoff teams

## Concurrency And Protection Mechanisms

The repo already uses optimistic locking for season state changes. That remains the correct default here because `Season` is still the aggregate root.

Recommended protections:

- `Season` version check for:
  - configuring playoff rules
  - generating playoff bracket
  - recording playoff leg results
  - advancing winners
  - finishing season

- database uniqueness constraints for:
  - one playoff ruleset per season
  - one standings snapshot per season playoff generation
  - one champion record per season if persisted separately

What race conditions matter most:

- two admins generating the playoff bracket at the same time
- concurrent updates to different legs that both try to resolve the same tie
- one admin changing rules while another starts playoffs
- one command finishing the season while another records a score

Minimal mechanism:

- optimistic locking on the `Season` aggregate version
- authoritative uniqueness constraints where separate tables are introduced

Do not split bracket advancement into ad hoc repository updates outside the season aggregate, or the invariants above will become fragile.

## Domain Methods Direction

The aggregate should expose explicit transitions such as:

- `ConfigurePlayoffRules(rules PlayoffRules) error`
- `CompleteRegularSeason(snapshot StandingsSnapshot) error`
- `GeneratePlayoffBracket() error`
- `RecordPlayoffLegScore(tieID string, legNumber int, score Score) error`
- `ResolveTie(tieID string, resolution TieResolutionInput) error`
- `AdvanceWinner(tieID string) error`
- `FinishWithChampion(teamID string) error`

Potentially better, some of these can be grouped so that the aggregate itself derives the next step:

- `RecordPlayoffLegScore(...)` may recompute aggregate score
- `ClosePlayoffTie(...)` may determine winner once enough information exists
- `AdvanceResolvedTie(...)` may place winner into the next round

The important point is that the domain methods should reflect business transitions and preserve invariants atomically at the aggregate level.

## Suggested Lifecycle

### Pre-Season

- Admin creates season
- Admin chooses playoff rules
- Backend validates rules

### Regular Season End

- Regular season finishes
- Final standings are frozen for qualification purposes
- Qualified teams are selected
- Seeds are assigned
- Playoff bracket is generated

### Playoffs In Progress

- Matches are scheduled round by round
- Tie winner is determined only when all required legs are complete
- Winners advance to the next round

### Season Finish

- Final tie winner becomes champion
- Season status becomes finished
- Champion and playoff outcome are stored in a season summary

## Validation Rules

The backend should reject invalid playoff configurations early.

Examples:

- qualifier count must be at least `2`
- qualifier count must be compatible with bracket shape
- final cannot have more than one tie
- two-leg rounds must define aggregate tie resolution
- third-place match requires semifinal losers to be known
- playoff rules cannot be changed once playoffs have started
- regular-season standings must be finalized before qualification is computed
- playoff bracket cannot be regenerated after playoff matches have started

## Data Modeling Direction

MVP-friendly direction:

- keep `Season` as the aggregate root
- add a `playoff_rules` configuration owned by the season
- store generated playoff bracket data separately from regular-season rounds
- represent playoff ties explicitly, with one or more legs per tie

Do not model playoffs as a separate season, and do not overload current `Round -> Matches` alone to represent multi-leg knockout logic. That structure is too weak for:

- aggregate scoring
- tie-level winner computation
- bracket advancement
- third-place logic

A playoff tie is a first-class concept and should exist in persistence and domain logic.

## Suggested Persistence Concepts

Likely tables or equivalents:

- `season_playoff_rules`
- `season_playoff_rounds`
- `season_playoff_ties`
- `season_playoff_legs`

Minimal information needed:

- season id
- round name and order
- seed numbers
- home and away teams per leg
- leg order
- aggregate scores
- winner team id
- advancement target
- tiebreak resolution used

## Read Models Needed

The frontend and operations side will likely need:

- playoff qualification summary
- generated bracket
- tie details with both legs
- aggregate score view
- champion summary
- season summary combining regular season and playoffs

## Admin Actions Needed

- configure playoff rules
- preview qualifiers and bracket before activation
- generate bracket
- schedule playoff matches
- record playoff match results
- resolve tie if manual decision is allowed
- publish champion

## Edge Cases To Plan For

- tied standings around the qualification line
- incomplete or cancelled regular-season matches before playoff generation
- odd number of qualifiers requiring byes
- team withdrawal after qualification
- leg cancellation or postponement
- neutral-site finals
- penalties recorded after a draw without changing regulation score

## Backend Design Notes

- Qualification should be computed from a finalized standings snapshot, not a mutable query that can drift later
- Playoff generation should be deterministic from `PlayoffRules + standings snapshot`
- Tie winner logic should live in domain code, not controllers
- Match-level score updates must not silently decide a playoff winner until the tie rules allow it
- Concurrency rules will matter once multiple playoff matches or tie updates occur at the same time
- `Season` should remain the aggregate root for playoff progression to preserve invariants cleanly

## Recommended MVP Slice

Build in this order:

1. Add season-level playoff rules configuration with validation
2. Support `top_n` qualification from final standings
3. Generate a fixed knockout bracket for `4` or `8` teams
4. Support one-leg rounds
5. Add two-leg semifinal support with aggregate winner computation
6. Store season champion and playoff summary

This sequence delivers value early while keeping the first backend iteration coherent.

## Open Decisions

Before implementation starts, the team should explicitly choose:

- whether MVP supports only `4` and `8` team playoffs or any even `N`
- whether finals can be two-leg
- whether away goals exist at all
- whether penalty winners need separate score fields
- whether admins can override qualifiers or seeds
- whether playoff bracket generation is automatic or manual

## Summary

The correct planning direction is a season-owned playoff rules model inside the same season, not a separate derived season, not a hardcoded post-season flow, and not a free-form scripting engine.

For MVP, treat it as a typed rules system that supports:

- top-`N` qualification
- standings-based seeding
- fixed knockout bracket
- one-leg and two-leg ties
- aggregate winner resolution
- champion determination

That is enough to support meaningful playoffs while keeping the design compatible with a broader multi-format season engine later.
