# Mobile App

This app is run from the monorepo root with Nx.

Browser development:

```bash
NX_DAEMON=false npx nx run mobile:web
```

Static web export:

```bash
NX_DAEMON=false npx nx run mobile:build-web
```

Tests:

```bash
NX_DAEMON=false npx nx test mobile
```

UI guidelines:

- `react-native-paper` is the default UI layer for text, buttons, cards, fields, search, surfaces, and overlays.
- Use app wrappers from `components/ui/` before reaching for raw Paper components in screen code.
- `StyledModal` owns `Portal` usage. Do not wrap modal callers in another `Portal`.
- Use `useAppTheme()` as the only theme API for product UI.
- Raw React Native primitives are for layout and low-level platform hooks, not for ad hoc button/text systems.

Use the root `README.md` for the full local setup flow.
