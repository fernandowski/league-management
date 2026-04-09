# Mobile UI Rules

## Default UI Layer

- `react-native-paper` is the default UI system for the mobile app.
- Prefer app wrappers in `components/ui/` before using raw Paper components in feature code.
- Raw React Native primitives are for layout, measurement, scrolling, and platform hooks, not for inventing parallel button/text/input systems.

## Approved Primitives

- Text: `AppText`
- Buttons/actions: `AppButton`
- Cards/surfaces: `AppCard`
- Overlays: `StyledModal` / `AppModal`
- Text fields: `AppTextField` and `ControlledTextInput`
- Selects: `AppSelect` / `Select`
- Screen shells: `AppScreen`

## Theme Rules

- Use `useAppTheme()` as the only theme API for product UI.
- Do not reintroduce alternate theme helper layers like `ThemedText`, `ThemedView`, or color maps outside the Paper theme.
- Put reusable visual tokens in `theme/theme.tsx`, not scattered inline across feature code.

## Modal Rules

- `StyledModal` owns the `Portal`.
- Do not wrap modal callers in another `Portal`.
- Do not create ad hoc `react-native-paper` `Modal` usage in feature files unless the shared modal abstraction is being extended.

## Interaction Rules

- Prefer Paper-backed interactions or app wrappers over `TouchableOpacity` for standard buttons, tabs, and action chips.
- If a raw `Pressable` is used, it should be for a specific interaction pattern that is not covered cleanly by the shared UI layer.
- Keep raw `View` usage focused on layout composition.

## Migration Direction

- When touching old mobile UI code, move it toward the shared wrapper layer instead of adding more direct Paper usage.
- Consolidate repeated card/modal/field styles into wrappers or shared variants.
- Remove dead starter/template UI rather than adapting it.

## Verification

- Run `npx nx lint web` after UI changes.
- If tests exist for the affected area, run `npx nx test web`.
- Treat lint warnings in touched files as cleanup work, not permanent debt.

## Device Support
- Always mobil-first then browser.
