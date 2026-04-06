# League Management

Nx monorepo with:

- `apps/api`: Go backend
- `apps/mobile`: frontend app, exposed in Nx as the `web` project

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
NX_DAEMON=false npx nx run api:migrate
```

These commands run migrations in Docker via `docker compose run migrate`.
`nx serve api` and `nx run api:dev` start the API container, but they do not run migrations automatically.

### 4. Start the backend

```bash
NX_DAEMON=false npx nx serve api
```

Backend runs on `http://localhost:8080`.
This command now starts the backend inside the Docker Compose `app` container.

### 5. Start the frontend for browser development

```bash
NX_DAEMON=false npx nx serve web
```

Frontend runs on `http://localhost:8081`.

## Useful Commands

```bash
NX_DAEMON=false npx nx test api
NX_DAEMON=false npx nx test web
NX_DAEMON=false npx nx run web:lint
NX_DAEMON=false npx nx build web
```

## Notes

- Use `NX_DAEMON=false` if Nx daemon issues appear locally.
- `apps/mobile` is an Expo Router app. In Nx it is now managed as the `web` project.
