# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Listinator is a shopping/to-do list web app. Go backend (Echo + GORM/SQLite), Vue 3 + TypeScript frontend built with Vite. The frontend build output is embedded into the Go binary via `//go:embed frontend/dist/*` (see `main.go`), so the shipped artifact is a single static binary.

## Commands

### Full build / run (from repo root)
```bash
make clean run   # clean artifacts, build frontend, build backend, run with dev config
make build        # build frontend then backend binary (./listinator)
make test         # build frontend (needed for the go:embed), then run Go and frontend tests
```
`make run` sets `LISTINATOR_SESSION_SECRET`, `LISTINATOR_ADMIN_PASSWORD`, and `LISTINATOR_DATABASE_DIR=.` for you.

### Backend (Go, root module)
```bash
go build ./...
go vet ./...
go build -o listinator .
go test ./...        # needs frontend/dist to exist, the main package embeds it
```
The only Go tests are `database/migration_test.go` — they run the migrations against a fresh
SQLite file in a temp dir and check the seeded entry types. Everything else is untested.

### Frontend (`frontend/`)
```bash
npm install
npm run dev          # vite dev server, proxies /api to http://localhost:8080
npm run build         # runs type-check (vue-tsc --build) and vite build in parallel
npm run type-check    # vue-tsc --build only
npm run test          # vitest run
npm run preview
```
No lint script and no ESLint config are set up in this repo.

Vitest is set up (`vitest.config.ts` reuses `vite.config.ts` for the `@/` alias, `environment: "node"`), and tests live in `src/**/__tests__/*.test.ts` — currently only for `src/i18n/`. There is no component or E2E test setup (no `@vue/test-utils`, jsdom, Playwright, Cypress), so UI changes still need to be verified manually by running the app (`make clean run` or `npm run dev`) and clicking through the flow in a browser.

### Database migrations
Migrations use [goose](https://github.com/pressly/goose) and live in `database/migrations/`. They are embedded (`//go:embed migrations/*.sql`) and run automatically on backend startup (`database/migration.go`). GORM auto-migration is intentionally **not** used.

To add a schema change, create a new migration and hand-write the `up` steps:
```bash
goose -dir database/migrations/ create do-something sql
```
Downgrades are not supported — only write the `Up` migration.

### Required environment variables (backend)
- `LISTINATOR_DATABASE_DIR` — directory for the SQLite file (required)
- `LISTINATOR_SESSION_SECRET` — session cookie signing secret (required)
- `LISTINATOR_ADMIN_PASSWORD` — if set, an `admin` user is created/updated with this password on every boot (`database/init.go`)
- `LISTINATOR_LOG_LEVEL` — `debug`|`info`|`warning`|`error` (default `info`)
- `LISTINATOR_LOG_TYPE` — `text`|`json` (default `text`)

### Docker
`Dockerfile` is a 3-stage build: Node stage builds the frontend, Go stage compiles the backend (`CGO_ENABLED=1`, needed for `mattn/go-sqlite3`) and embeds `frontend/dist`, final stage is `debian:bookworm-slim` with just the binary. The Go base image tag must satisfy the `go` directive version in `go.mod` — bumping the toolchain requirement in `go.mod` (e.g. via dependency updates) requires bumping the `FROM golang:...` tag in the Dockerfile too, or the Docker build breaks with a `GOTOOLCHAIN=local` version error.

CI (`.github/workflows/docker-build.yml`) has a `test` job (frontend `npm test`, then `npm run build` so `go test ./...` can compile the embed) that gates the `build-and-push` job, which builds and pushes the Docker image on every push to any branch (tagged `:<sha>`, plus `:latest` on `main` and `:<tag>` on tag pushes). There is no lint job.

## Architecture

### Backend layout
- `main.go` — wires everything: logger init, DB init, Echo instance, session middleware (`gorilla/sessions` cookie store via `echo-contrib/session`), mounts API routes under `/api/v1`, and serves the embedded frontend at `/`.
- `api/v1/server/` — HTTP layer. `server.go` defines the `server` struct (holds `*gorm.DB` and an in-memory pub/sub for entry events) and `SetupRoutes`. Each resource (entries, lists, types, session) has its own file with closures returning `echo.HandlerFunc`; request structs are defined inline per-handler and bound via Echo's `c.Bind`.
- `database/` — GORM models (`models.go`), SQLite init + admin-user bootstrap (`init.go`), goose migration runner (`migration.go`), and a custom `slog`-backed GORM logger (`logger.go`).
- `pubsub/` — small generic in-memory pub/sub (`PubSub[K, T]`), used to fan out entry create/update/delete events to Server-Sent-Event subscribers per list.
- `logger/` — process-wide `slog` setup from env vars.

### Data model (`database/models.go`)
All models embed `Model` (UUID primary key generated in `BeforeCreate`, timestamps, soft-delete). Hierarchy: `User` (auth), `List` has many `Entry`, `Entry` belongs to a `Type` (category/priority, used for sorting/coloring). Entries without an explicit `TypeID` default to a hardcoded "misc" type UUID.

### Real-time updates (SSE)
Clients poll no; instead `GET /api/v1/entries/events?ListID=...` (`entryGetEvents` in `entry.go`) opens a Server-Sent-Events stream. Every entry create/update/delete handler publishes to `server.entryPubSub` keyed by `ListID`; the SSE handler subscribes for the connection's lifetime, forwards events, and sends a ping every 5s to keep the connection alive. This is how the frontend gets live updates without polling.

### Auth
Session-based, not JWT: `POST /api/v1/session` checks bcrypt password hash and stores the user's UUID in a signed cookie session (`gorilla/sessions`). `sessionMiddleware` in `session.go` resolves the UUID back to a `database.User` and puts it in the Echo context (`userKey`). Only the `GET /api/v1/session` route currently uses this middleware — other routes are not currently gated by it.

### Frontend layout (`frontend/src/`)
Vue 3 + `<script setup>` SPA using hash-based routing (`createWebHashHistory`, see `router.ts` — this is why the app can be served from a static file server / embedded FS without server-side routing config). `src/api/api.ts` is a thin fetch wrapper (`apiFetchJSON`) around the `/api/v1/*` endpoints; `vite.config.ts` proxies `/api` to `localhost:8080` in dev. Pages live in `src/Pages/`, shared UI in `src/Components/`, layout shell in `src/Layouts/`.

### Internationalization (`frontend/src/i18n/`)
English and German, hand-rolled instead of `vue-i18n` — the long-term intention is to be able to drop Vue, so `core.ts` (locale resolution, catalog lookup, `{placeholder}` formatting, fallback chain) is pure TypeScript with no framework or browser access, and `index.ts` is the only part touching Vue (`locale` ref), `localStorage`, and `document.documentElement.lang`.

Every user-visible string goes into `locales/en.ts` — it is the source of truth and `locales/de.ts` is typed as `Messages = typeof en`, so a missing key fails `npm run type-check`. Components use a plain function import, not a composable: `import { t } from "@/i18n"` and then `t("list.search")`. Anything holding a translated string outside the template must be a `computed` (see `contextmenuActions` in `Pages/List.vue`), otherwise it freezes at the language active during the first render.

Entry type names are **not** translated in the backend: the `types` table has a `slug` column (seeded in the `add_type_slug` migration) and `translateType()` looks up `types.<slug>` in the catalog, falling back to the `Name` column for unknown slugs. A new type therefore needs both a migration with a slug and an entry in every catalog — `src/i18n/__tests__/slugs.test.ts` reads the migration files and fails if the two sides drift apart.

### Styling
Global design tokens (color variables, `--border-radius`) and base element styles (`header`, `input`, `body`, etc.) are defined unscoped in `Layouts/DefaultLayout.vue`'s `<style>` block and apply app-wide. Individual components then override these locally in their own scoped `<style>` — e.g. `Components/EntryItem.vue` overrides the base `input` style with a green accent border. When matching an existing look, check `DefaultLayout.vue` first for the shared base, then check for component-level scoped overrides before assuming a style is component-local.

### Deployment reference
`infra/docker-compose.yaml` and `infra/conf/Caddyfile` show the intended production setup (Caddy in front of the container) — useful as a reference when reasoning about env/config, not something built by this repo's CI.
