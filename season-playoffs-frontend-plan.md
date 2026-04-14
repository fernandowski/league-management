# Season Playoffs Frontend Plan

This document plans the frontend work needed to support season playoffs in `apps/mobile`.

It is based on the current season UI, which already has:

- season details view
- season status and basic metrics
- standings tab
- match-up management tab
- season planning and start actions

Relevant current files:

- [apps/mobile/app/dashboard/seasons/[id].tsx](/home/fernando/ideaProjects/league-management/apps/mobile/app/dashboard/seasons/[id].tsx:1)
- [apps/mobile/components/Seasons/SeasonManagement.tsx](/home/fernando/ideaProjects/league-management/apps/mobile/components/Seasons/SeasonManagement.tsx:1)
- [apps/mobile/components/Seasons/SeasonDetails.tsx](/home/fernando/ideaProjects/league-management/apps/mobile/components/Seasons/SeasonDetails.tsx:1)
- [apps/mobile/components/Seasons/SeasonInformation.tsx](/home/fernando/ideaProjects/league-management/apps/mobile/components/Seasons/SeasonInformation.tsx:1)
- [apps/mobile/components/Seasons/SeasonStanding.tsx](/home/fernando/ideaProjects/league-management/apps/mobile/components/Seasons/SeasonStanding.tsx:1)
- [apps/mobile/components/Seasons/SeasonMatchUpManagement.tsx](/home/fernando/ideaProjects/league-management/apps/mobile/components/Seasons/SeasonMatchUpManagement.tsx:1)
- [apps/mobile/hooks/useData.ts](/home/fernando/ideaProjects/league-management/apps/mobile/hooks/useData.ts:1)

## Frontend Goals

- let an admin understand whether a season has playoffs enabled
- let an admin configure playoff rules before playoffs start
- show when the regular season is complete and playoffs are ready
- let an admin preview qualifiers and seeds before bracket generation
- show the playoff bracket clearly
- let an admin record playoff leg scores and see aggregate results
- show the final champion in the season summary

## Current UI Shape

Today the season experience is split into two patterns:

### Season Detail Card Pattern

Used in `SeasonDetails.tsx` and `SeasonInformation.tsx`.

Current responsibilities:

- season title and status
- rounds count
- start season action
- plan season action
- embedded schedule and standings

### Season Tabs Pattern

Used in `SeasonManagement.tsx`.

Current tabs:

- `Standings`
- `Match Ups`

This existing tabbed structure is the right place to grow the playoff UX.

## Recommended UX Structure

Do not create a separate standalone playoff screen first. Keep playoffs under the season detail flow.

Recommended season navigation:

- `Overview`
- `Standings`
- `Regular Season`
- `Playoffs`

For MVP, the simplest adjustment is:

- keep the current season page
- add a `Playoffs` tab
- move top-level season summary into a slightly richer overview area

## Proposed Screen Responsibilities

### 1. Season Overview

Purpose:

- summarize current season state
- surface the next available admin action
- show whether playoffs are configured

Content:

- season name
- lifecycle status
- season phase
- rounds count
- playoff enabled yes/no
- playoff format summary
- champion summary when finished
- next primary action

Current base:

- `SeasonInformation.tsx`

Recommended evolution:

- expand `SeasonInformation` into a proper overview summary card
- add a secondary card for phase-specific actions

### 2. Standings

Purpose:

- show regular-season table
- indicate whether standings are still live or frozen for playoffs
- eventually show qualification cut line

Current base:

- `SeasonStanding.tsx`

Changes needed:

- support qualification markers such as seed number or “qualified”
- support a “frozen for playoffs” label
- support richer standings fields later without redesigning the table again

### 3. Regular Season

Purpose:

- manage regular-season rounds and scores
- stop being the only match-management surface once playoffs exist

Current base:

- `SeasonMatchUpManagement.tsx`

Changes needed:

- label this explicitly as regular-season match management
- disable editing once regular season is completed
- show a clear transition state when playoffs are ready

### 4. Playoffs

Purpose:

- configure rules
- preview qualifiers
- generate bracket
- manage playoff scores
- track aggregate results and champion

This will be the main new frontend surface.

## Proposed Playoffs Tab Layout

One large card with sections is better than several disconnected cards.

Recommended order inside the tab:

### Section A: Playoff Summary

- playoffs enabled yes/no
- qualifier count
- rounds format summary
- current playoff status
- champion, when available

### Section B: Configuration

Visible when:

- playoffs not yet started

Fields:

- qualifier count
- semifinal legs
- final legs if supported
- third-place toggle if supported
- tie resolution summary

Actions:

- save rules
- discard changes

### Section C: Qualification Preview

Visible when:

- regular season completed
- standings snapshot available or preview endpoint available

Content:

- seed number
- team name
- points
- relevant ranking metrics
- cut line marker

Action:

- generate bracket

### Section D: Bracket

Visible when:

- bracket exists

Content:

- round columns or stacked round sections on smaller screens
- each tie shown as a card
- each leg listed inside the tie
- aggregate score summary
- winner badge when resolved

### Section E: Match Actions

Visible when:

- the user can record playoff scores

Content:

- score-edit action on each playable leg
- disabled state on completed or unresolved/future ties

## Recommended Component Breakdown

Reuse the current style of small season-focused components rather than introducing a huge screen file.

Suggested new components:

- `PlayoffSummaryCard`
- `PlayoffRulesForm`
- `PlayoffQualificationPreview`
- `PlayoffBracketView`
- `PlayoffRoundSection`
- `PlayoffTieCard`
- `PlayoffLegRow`
- `ResolvePlayoffTieModal` if manual resolution is needed

Suggested existing components to extend:

- `SeasonInformation.tsx`
- `SeasonManagement.tsx`
- `ManageMatchScoreModal.tsx`
- `TableList.tsx`
- `Tabs.tsx`

## Reuse Strategy

### Reuse `Tabs`

Add a `Playoffs` tab rather than inventing new internal navigation.

### Reuse score modal patterns

`ManageMatchScoreModal.tsx` already establishes the score editing interaction. It should be reused or lightly generalized for playoff leg score updates.

### Reuse cards and tables

- `AppCard`
- `AppButton`
- `AppText`
- `TableList`

These are enough for the first playoff UI.

## Data Contract Changes Needed From Backend

The frontend should not build playoffs from ad hoc reuse of current season endpoints. It needs explicit playoff-oriented payloads.

Recommended query payloads:

### Season Summary

Needed additions:

- `phase`
- `playoff_rules_summary`
- `playoff_status`
- `champion`

### Standings

Needed additions:

- `rank`
- `qualified`
- `seed` if qualification preview is frozen
- `snapshot_status`

### Playoff Qualification Preview

Suggested payload:

- `season_id`
- `qualifier_count`
- `qualified_teams[]`
  - `team_id`
  - `team_name`
  - `seed`
  - `rank`
  - `points`

### Playoff Bracket

Suggested payload:

- `season_id`
- `status`
- `winner_team_id`
- `rounds[]`
  - `name`
  - `order`
  - `ties[]`
    - `id`
    - `home_team`
    - `away_team`
    - `home_seed`
    - `away_seed`
    - `status`
    - `aggregate_score`
    - `winner_team_id`
    - `legs[]`
      - `match_id`
      - `leg_number`
      - `home_team`
      - `away_team`
      - `home_score`
      - `away_score`
      - `status`

### Playoff Rules

Suggested payload:

- `qualifier_count`
- `third_place_match`
- `rounds[]`
  - `name`
  - `legs`
  - `tied_aggregate_resolution`

## Hooks And Client Types To Add

Current `useData.ts` types are too narrow and regular-season-specific.

Add new types for:

- `SeasonSummaryResponse`
- `PlayoffRulesResponse`
- `PlayoffQualificationPreviewResponse`
- `PlayoffBracketResponse`
- `PlayoffTie`
- `PlayoffLeg`

The frontend should not keep playoff payloads in `Record<string, any>` style structures.

## Proposed Component and Route Changes

### `SeasonManagement.tsx`

Change from:

- `Standings`
- `Match Ups`

To something like:

- `Overview`
- `Standings`
- `Regular Season`
- `Playoffs`

If keeping the current detail card on another route, then at minimum:

- `Standings`
- `Regular Season`
- `Playoffs`

### `SeasonInformation.tsx`

Add:

- season phase
- playoff enabled label
- playoffs ready state
- champion summary

### `SeasonStanding.tsx`

Add:

- rank column
- qualification status display
- frozen standings messaging

### `SeasonMatchUpManagement.tsx`

Add:

- explicit regular-season wording
- disabled completion/edit state after regular season ends
- handoff banner when playoffs become available

### New `SeasonPlayoffs.tsx`

This should coordinate:

- rules fetch/save
- qualification preview
- bracket fetch
- bracket actions

## UX States To Design Explicitly

The tab must not assume one happy path.

### State 1: No Playoff Rules

Show:

- explanatory summary
- rules form
- save action

### State 2: Rules Configured, Regular Season In Progress

Show:

- rules summary
- note that bracket is unavailable until regular season completion

### State 3: Regular Season Complete, Bracket Not Generated

Show:

- qualified teams preview
- generate bracket action

### State 4: Bracket Generated, Playoffs Not Started

Show:

- bracket
- start playoffs action if backend exposes it

### State 5: Playoffs In Progress

Show:

- bracket with leg status
- score entry actions
- aggregate summaries

### State 6: Season Finished

Show:

- final bracket
- champion card
- regular season and playoff summary

## Interaction Plan

### Configure Rules

- simple inline form in the playoffs tab
- optimistic save feedback
- read-only once playoffs begin

### Generate Bracket

- high-emphasis action near qualification preview
- confirmation modal optional, not required for MVP

### Record Playoff Leg Score

- reuse modal interaction from regular-season score entry
- refresh bracket after save

### Resolve Tie

Only needed if MVP supports manual resolution.

If manual resolution is not in MVP, do not add this UI yet.

## Responsive Layout Plan

The current app already supports stacked versus side-by-side layouts in some season screens.

For playoffs:

- on larger screens:
  - use one summary/config card and one bracket card
  - bracket may display rounds side by side if simple enough

- on smaller screens:
  - stack sections vertically
  - render bracket as round sections, not a dense horizontal bracket diagram

Do not build a complex desktop-only bracket visualization first. A stacked mobile-first bracket is enough for MVP.

## Frontend Delivery Order

### Phase 1: Data Contracts And Types

- add client types for season phase and playoff payloads
- update season summary payload usage

### Phase 2: Season Overview Updates

- add phase and playoff summary to season overview
- add champion support

### Phase 3: Playoffs Tab Shell

- add new tab
- add state-driven empty/configured/ready/in-progress/finished views

### Phase 4: Rules Configuration

- add playoff rules form
- save and refresh rules

### Phase 5: Qualification Preview

- show qualified teams and seeds
- add bracket generation action

### Phase 6: Bracket View

- render rounds, ties, and legs
- show aggregate scores and winners

### Phase 7: Playoff Match Management

- support score editing on playoff legs
- refresh bracket state after updates

### Phase 8: Season Completion UX

- champion card
- finished-state messaging

## MVP Cut Recommendation

If frontend scope needs to stay tight, the MVP UI should include only:

- playoffs summary in season overview
- playoffs tab
- rules form
- qualification preview
- bracket view
- score entry for playoff legs
- champion display

Do not include in MVP:

- animated bracket lines
- drag-and-drop seeds
- manual tie override UI unless backend requires it
- separate public-facing playoff pages

## Main Frontend Risks

- trying to force playoffs into the current regular-season match tab
- overbuilding a fancy bracket visualization before core flows work
- using weak client types and then spreading `any` through the season UI
- letting regular season and playoff actions appear active at the same time

## Summary

The cleanest frontend path is to extend the existing season experience, not replace it.

Recommended shape:

- richer season overview
- current standings preserved and upgraded
- regular-season match management stays separate
- new dedicated `Playoffs` tab for rules, qualification preview, bracket, and playoff scoring

That fits the current component structure and will scale cleanly once the backend playoff model lands.
