package internal

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres to migration, err: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect postgres to migration, err: %w", err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run postgres migration, err: %w", err)
	}
	return nil
}
