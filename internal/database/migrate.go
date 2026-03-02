package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func MigrateUp(dbURL string, dir string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}
	return nil
}

func MigrateDown(dbURL string, dir string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Down(db, dir); err != nil {
		return fmt.Errorf("goose down error: %w", err)
	}
	return nil
}
