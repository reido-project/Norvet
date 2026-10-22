package provider

import (
	"api/persistence/config"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPostgres(cfg config.Config) (*sql.DB, string, error) {
	if cfg.Postgres == nil {
		return nil, "", fmt.Errorf("postgres config is missing")
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	q := u.Query()
	q.Set("sslmode", cfg.Postgres.SSLMode)
	q.Set("TimeZone", cfg.Postgres.Timezone)
	u.RawQuery = q.Encode()

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, "", fmt.Errorf("open PostgreSQL: %w", err)
	}

	return db, "postgres", nil
}
