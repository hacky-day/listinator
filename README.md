# Listinator

A simple web-based list management application perfect for shopping lists, to-do lists, or any kind of item tracking.
Backends is built with Go and a modern frontend with Typescript and Vue embeded in the Go Binary for easy deployment.

---

## Features

- Go Backend with persistent storage via SQLite
- Web interface for managing lists

## Setup (Production)

For production, use the pre-built Container Images. The
[docker-compose.yaml](./infra/docker-compose.yaml) can be used as a reference for
production deployment.

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE)
file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
