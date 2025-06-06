package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/Axel791/auth/internal/config"
)

// ConnectDB - подключение к базе данных, применение миграций
func ConnectDB(databaseDSN string, cfg *config.Config) (*sqlx.DB, error) {
	if databaseDSN != "" {
		db, err := sqlx.Connect("postgres", databaseDSN)
		if err != nil {
			return nil, err
		}

		err = AppleMigration(db, cfg)
		if err != nil {
			_ = db.Close()
			return nil, err
		}

		return db, nil
	}
	return nil, nil
}
