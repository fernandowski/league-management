# Features Backlog

This document captures feature ideas that fit the current product shape of the league management project.

The app already supports:

- Organizations
- Teams
- Leagues
- League membership management
- Seasons
- Round-robin schedule generation
- Season start
- Match score updates
- Standings

## Prioritization Principles

Prioritize features that:

- Extend existing domain concepts already present in the codebase
- Increase day-to-day usefulness for league operators
- Open the product to more user roles than just the organization owner
- Add operational depth before introducing large new domains

## Quick Wins

### 1. Match Scheduling Details

Add real fixture logistics on top of generated rounds.

Scope:

- Match date and kickoff time
- Venue or field assignment
- Postponed and cancelled match states
- Manual rescheduling
- Upcoming matches view

Why it matters:

- Makes schedules usable in real life
- Builds directly on the current season and match model
- Unlocks reminders and calendar features later

### 2. Full Season Lifecycle

Expand the season workflow beyond pending, planned, and in-progress.

Scope:

- Pause season
- Resume season
- Finish season
- Prevent score changes after finish
- Show champion and final table

Why it matters:

- The domain already includes paused and finished states
- Closes the loop for a season from creation to completion

### 3. Referee Assignment

Finish the referee workflow already hinted at in the domain and persistence layer.

Scope:

- Create and manage referees
- Assign a referee to a match
- Show referee on fixture details
- Restrict referee score updates to assigned matches

Why it matters:

- The codebase already includes referee-related models and fields
- Adds a realistic operational workflow without requiring a major redesign

### 4. Better Standings Rules

Make standings logic configurable and competition-ready.

Scope:

- Goal difference
- Goals against
- Head-to-head tiebreakers
- Configurable points systems
- Manual standings override for admin correction

Why it matters:

- Different leagues use different ranking rules
- Makes the product more reusable across formats

## Medium Projects

### 5. Team Manager Access

Move beyond the current owner-only administration pattern.

Scope:

- Team manager role
- Team-scoped permissions
- Team dashboard for schedule and results
- Team season summary view

Why it matters:

- Expands the number of users who benefit from the product
- Turns the app into a participant tool, not just an admin console

### 6. Organization and League Roles

Introduce role-based access control across the product.

Scope:

- Organization admin
- League manager
- Team manager
- Referee
- Permission checks by route and service

Why it matters:

- The current app is heavily centered on the organization owner
- Real leagues usually distribute work across multiple people

### 7. Notifications and Reminders

Add communication features around key league events.

Scope:

- Match reminders
- Season start notifications
- Score submission reminders
- Membership invite status notifications
- In-app notification center

Why it matters:

- Improves operational follow-through
- Builds on scheduling and role features

### 8. Public Fixtures and Standings Pages

Expose read-only league data publicly.

Scope:

- Public standings page
- Public fixtures and results
- Public league summary page
- Shareable links

Why it matters:

- Increases product visibility and usefulness
- Does not require all users to authenticate

## Bigger Bets

### 9. Playoffs and Knockout Stages

Extend seasons beyond round-robin regular play.

Scope:

- Seed teams from standings
- Generate semifinals and finals
- Bracket display
- Third-place match support
- Champion history

Why it matters:

- Many leagues end with a playoff stage
- This is a natural next format after standings are stable

### 10. Rosters and Player Management

Model the people inside teams, not just teams themselves.

Scope:

- Player profiles
- Roster management
- Jersey numbers
- Player eligibility
- Transfers or roster locks

Why it matters:

- Deepens the team model significantly
- Opens the door to player stats and disciplinary systems

### 11. Player Statistics and Discipline

Track match-level player outcomes.

Scope:

- Goals
- Assists
- Cards
- Suspensions
- Top scorer and leaderboard views

Why it matters:

- Adds engagement for participants
- Creates more value from each recorded match

### 12. Multi-Format Season Engine

Generalize the season model to support more competition structures.

Scope:

- Double round robin
- Groups plus knockout
- Custom number of rounds
- Manual fixture templates

Why it matters:

- Makes the product usable for more league types
- Requires careful domain and persistence design

## Recommended Delivery Order

### Phase 1: Operational Depth

- Match scheduling details
- Full season lifecycle
- Referee assignment

### Phase 2: More Users

- Team manager access
- Organization and league roles
- Notifications and reminders

### Phase 3: Audience Growth

- Public fixtures and standings pages

### Phase 4: Competition Depth

- Playoffs and knockout stages
- Rosters and player management
- Player statistics and discipline
- Multi-format season engine

## Highest-Leverage Next Step

If only one feature should be built next, start with match scheduling details.

Reason:

- It makes the existing schedule generation actually useful for real operations
- It sets up future work for reminders, referee assignment, and public fixtures
- It fits the current domain model without forcing a large architectural change

## Possible Follow-Up Planning Docs

Future planning documents that could be useful:

- `roadmap.md` for sequencing by milestone
- `roles-and-permissions.md` for access-control design
- `season-formats.md` for competition structure options
- `notifications.md` for reminder and messaging workflows
