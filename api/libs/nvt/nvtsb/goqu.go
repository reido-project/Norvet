package sqlbuilder

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"

	"api/libs/nvt"
)

// ConnectGoqu connects to the database configured in the environment (see
// nvtsql.LoadConfig) and returns a goqu query builder for it. nvt.Init must
// have been called first.
//
// Close the connection with CloseGoqu when the application shuts down.
func ConnectGoqu() (*goqu.Database, error) {
	return connect(WrapGoqu)
}

// WrapGoqu returns a goqu query builder for an existing *sql.DB of the given
// dialect. Use it when you already hold the handle, for example one from
// nvtsql.Open. The builder shares db; closing either closes the pool.
func WrapGoqu(db *sql.DB, dialect nvt.Dialect) (*goqu.Database, error) {
	name, err := goquDialect(dialect)
	if err != nil {
		return nil, err
	}
	return goqu.Dialect(name).DB(db), nil
}

// CloseGoqu closes the connection pool behind a builder created by
// ConnectGoqu or WrapGoqu.
func CloseGoqu(db *goqu.Database) error {
	sqlDB, ok := db.Db.(*sql.DB)
	if !ok {
		return fmt.Errorf("sqlbuilder: goqu database is not backed by *sql.DB (%T)", db.Db)
	}
	return sqlDB.Close()
}

// goquDialect maps a Norvet dialect to the name of its goqu SQL dialect. goqu
// silently falls back to a generic dialect for unknown names, so the mapping
// is checked here instead.
func goquDialect(dialect nvt.Dialect) (string, error) {
	switch dialect {
	case nvt.MySQL, nvt.MariaDB:
		return "mysql", nil
	case nvt.Postgres:
		return "postgres", nil
	case nvt.SQLite:
		return "sqlite3", nil
	}
	return "", fmt.Errorf("sqlbuilder: goqu does not support dialect %q", dialect)
}
