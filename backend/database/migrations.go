package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func RunMigrations(dbconn *sql.DB) error {
	// Настройка goose для использования postgres
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Выполнение миграций
	if err := goose.Up(dbconn, "database/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
