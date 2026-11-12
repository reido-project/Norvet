// Package postgres registers the PostgreSQL connector (pgx) with nvtsql.
// Import it for its side effect:
//
//	import _ "api/libs/nvt/nvtsql/postgres"
package postgres

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"api/libs/nvt"
	"api/libs/nvt/nvtsql"
)

func init() {
	nvtsql.Register(nvt.Postgres, open)
}

func open(cfg nvtsql.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("nvtsql: open postgres: %w", err)
	}
	return db, nil
}

// dsn builds a postgres:// URL. The session TimeZone is only set for a named
// location; "Local" is not a name PostgreSQL understands.
func dsn(cfg nvtsql.Config) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   "/" + cfg.Name,
	}

	q := u.Query()
	if cfg.Postgres.SSLMode != "" {
		q.Set("sslmode", cfg.Postgres.SSLMode)
	}
	if loc := cfg.Location; loc != nil && loc != time.Local {
		q.Set("TimeZone", loc.String())
	}
	u.RawQuery = q.Encode()

	return u.String()
}
