---
name: frontend-mobile-league
description: Use when working on the mobile frontend in this repo, especially under `apps/mobile`, for UI features, layout changes, component reuse, styling cleanup, interaction refinements, and frontend reviews. Trigger for screens, tabs, cards, tables, modals, buttons, responsive layout composition, and cases where existing mobile components should be reused instead of rebuilt.
---

# Frontend Mobile League

## Overview

Use this skill for frontend work in `apps/mobile`. Favor reuse of existing components and patterns over inventing new UI primitives. Keep changes intentional, compact, and consistent with the established mobile design system.

Start by reading `apps/mobile/AGENTS.md` if present, then inspect the target screen and nearby shared UI components before editing.

## Workflow

1. Find the existing route, screen, and shared components involved.
2. Reuse or extend an existing component before creating a new one.
3. Choose the smallest layout change that solves the request on desktop and mobile.
4. Promote repeated button or card treatments into shared UI APIs instead of duplicating styles.
5. Remove filler copy and redundant wrappers when they do not improve usability.
6. Run the narrowest useful verification, usually `npx expo lint` from `apps/mobile`.

## Repo Patterns

### Reuse First

- Search `apps/mobile/components` for an existing card, table, tab, modal, or stats component before building a new one.
- If a component is almost correct, add small extension points such as optional style props instead of cloning it.
- Keep behavior in the original shared component when the same UI can be needed in multiple screens.

### Layout Composition

- Prefer composing screens from `AppCard`, `AppText`, `AppButton`, existing tables, and existing feature components.
- On large screens, use a two-column top section when it improves scanability.
- On smaller screens, stack the same content vertically rather than creating a separate mobile-only structure.
- When combining related content, prefer one larger card with clear sections over multiple nested wrappers.
- Remove redundant inner cards when one parent card already defines the surface.

### Content Density

- Keep instructional copy short and only keep text that changes user decisions.
- Remove helper text that merely explains obvious controls or repeats current system behavior.
- Put primary actions close to the content they affect.
- When two controls belong together, such as pagination and a round-completion action, keep them on the same row when space allows.

### Shared UI Rules

- Use `AppButton` instead of raw Paper buttons.
- Prefer semantic button variants over one-off styling.
- Use `submit` for save, create, start, plan, and complete actions.
- Use `secondary` for cancel or close actions.
- Use `destructive` for remove or delete actions.
- If a new visual treatment appears more than once, add it to the shared component API instead of restyling locally.

### Styling Rules

- Preserve the established visual language in `apps/mobile/theme/theme.tsx`.
- Prefer theme colors over hard-coded values.
- Keep spacing and radii consistent with nearby code.
- Add style props to shared components only when they improve reuse without making the API vague.
- Avoid introducing decorative UI that does not support the task.

## Common File Targets

- Routes: `apps/mobile/app/**/*`
- Feature screens: `apps/mobile/components/Seasons/*`, `League/*`, `Overview/*`, `Teams/*`
- Shared UI: `apps/mobile/components/ui/*`
- Shared layout helpers: `apps/mobile/components/Layout/*`, `TableList/*`
- Theme tokens: `apps/mobile/theme/theme.tsx`

## Review Checklist

- An existing component was reused or extended before creating a new one.
- The layout works on both larger and smaller screens.
- Extra helper copy, wrappers, and duplicate actions were removed where possible.
- Shared primitives were updated when a pattern became reusable.
- Button semantics match the action intent.
- Lint was run if the change touched app code.
