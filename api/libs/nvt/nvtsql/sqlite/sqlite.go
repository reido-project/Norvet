// Package sqlite connects to SQLite databases and registers the SQLite
// connector with nvtsql.
//
// It is deliberately small: it needs only database/sql and the SQLite driver,
// so code that just wants a SQLite file (such as nvtlog) does not pull in a
// query builder or the other database drivers.
//
// The driver is github.com/mattn/go-sqlite3, which requires cgo.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"api/libs/nvt"
	"api/libs/nvt/nvtsql"
)

func init() {
	nvtsql.Register(nvt.SQLite, func(cfg nvtsql.Config) (*sql.DB, error) {
		return Open(cfg.Name)
	})
}

// Open opens the SQLite database at path (":memory:" for a private in-memory
// database) and verifies it. The handle is limited to one connection, which
// serializes writers and avoids "database is locked" errors. The caller must
// Close it.
func Open(path string) (*sql.DB, error) {
	if path == "" {
		return nil, errors.New("nvtsql: sqlite path is empty")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("nvtsql: open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nvtsql: ping sqlite: %w", err)
	}
	return db, nil
}
