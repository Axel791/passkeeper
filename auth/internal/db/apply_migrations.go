package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"github.com/Axel791/auth/internal/config"
)

// AppleMigration - Применение миграций
func AppleMigration(dbConn *sqlx.DB, cfg *config.Config) error {
	return goose.RunContext(context.Background(), "up", dbConn.DB, cfg.MigrationsPath)
}
