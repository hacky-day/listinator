package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("unable to set sqlite3 dialect in goose, %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("unable to migrate, %w", err)
	}

	return nil
}
