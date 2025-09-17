# Core

Core backend of the listinator with embedded frontend.

---

## Setup (Production)

For production, use the pre-built Container Images.

## Setup (Development)

### Requirements

- Go 1.24.3 or later
- Node.js 18+ and npm
- `make`

### Quick Start

Simply run:

```bash
make clean run
```

This will:

1. Clean old artifacts
2. Build the frontend
3. Build the backend and embed the frontend
4. Start the backend with proper Configuration for Development

## Database Migrations

Database migrations are done with [goose](https://github.com/pressly/goose).

This means this project does not use [Gorms
Automigration](https://gorm.io/docs/migration.html#Auto-Migration) anymore, but
on each database change a new migration needs to be created with

```bash
goose -dir database/migrations/ create do-something sql
```

and the migration steps need to be created manually.

We do not support downgrades, so writing the steps for `goose up` should be enough.

Migration scripts are automatically called on boot.

## Configuration

The application uses the following environment variables:

- `LISTINATOR_DATABASE_DIR` - Directory where the SQLite database file will be
  stored (required)
- `LISTINATOR_SESSION_SECRET` - Secret key used for session management
  (required)
- `LISTINATOR_ADMIN_PASSWORD` - Password for admin access (required)
- `LISTINATOR_LOG_LEVEL` - Log level for application logging. Options: `debug`,
  `info`, `warning`, `error`. Defaults to `info`
- `LISTINATOR_LOG_TYPE` - Log output format. Options: `text`, `json`. Defaults
  to `text`
- `LISTINATOR_TYPIFIER_URL` - URL of the typifier. If not set, the listinator
  is running without this microservice
