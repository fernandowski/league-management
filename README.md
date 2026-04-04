# League Management

Nx monorepo with:

- `apps/api`: Go backend
- `apps/mobile`: Expo app used for browser development and mobile targets

## Local Development

### 1. Install dependencies

```bash
npm install
```

### 2. Start Postgres

```bash
docker compose up -d db
```

### 3. Run database migrations

```bash
npm run migrate:up
```

### 4. Start the backend

```bash
NX_DAEMON=false npx nx serve api
```

Backend runs on `http://localhost:8080`.

### 5. Start the frontend for browser development

```bash
NX_DAEMON=false npx nx run mobile:web
```

Frontend runs on `http://localhost:8081`.

## Useful Commands

```bash
NX_DAEMON=false npx nx test api
NX_DAEMON=false npx nx test mobile
NX_DAEMON=false npx nx lint mobile
NX_DAEMON=false npx nx run mobile:build-web
```

## Notes

- Use `NX_DAEMON=false` if Nx daemon issues appear locally.
- `apps/mobile` is an Expo Router app. For browser work, use the `mobile:web` target.
