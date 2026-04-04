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

Use the root `README.md` for the full local setup flow.
