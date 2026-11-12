package sqlbuilder

import (
	"database/sql"

	"api/libs/nvt"
	"api/libs/nvt/nvtsql"
	_ "api/libs/nvt/nvtsql/drivers"
)

// wrapFunc adapts an open *sql.DB to a query builder of type T.
type wrapFunc[T any] func(db *sql.DB, dialect nvt.Dialect) (T, error)

// connect opens the database from the environment (see nvtsql.LoadConfig) and
// wraps it with wrap. If wrapping fails the database is closed, so callers
// never leak a connection pool.
//
// Every Connect* function is this call with its own wrap function; supporting
// another query builder never touches existing code.
func connect[T any](wrap wrapFunc[T]) (T, error) {
	var zero T

	cfg, err := nvtsql.LoadConfig()
	if err != nil {
		return zero, err
	}

	db, err := nvtsql.Open(cfg)
	if err != nil {
		return zero, err
	}

	builder, err := wrap(db, cfg.Dialect)
	if err != nil {
		_ = db.Close()
		return zero, err
	}
	return builder, nil
}
