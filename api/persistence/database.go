package persistence

import (
	"api/persistence/config"
	"api/persistence/provider"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func Connect() (*goqu.Database, error) {
	cfg := config.LoadConfig()

	sqlDB, dialect, err := openDatabase(cfg)
	if err != nil {
		return nil, err
	}

	var success bool
	defer func() {
		if !success {
			_ = sqlDB.Close()
		}
	}()

	configurePool(sqlDB, cfg)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	success = true
	return goqu.Dialect(dialect).DB(sqlDB), nil
}

func openDatabase(cfg config.Config) (*sql.DB, string, error) {
	switch cfg.Driver {
	case "mysql", "mariadb":
		return provider.OpenMySQL(cfg)

	case "postgres", "postgresql":
		return provider.OpenPostgres(cfg)

	default:
		return nil, "", fmt.Errorf("unsupported DB_DRIVER: %q", cfg.Driver)
	}
}

func configurePool(db *sql.DB, cfg config.Config) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
}
