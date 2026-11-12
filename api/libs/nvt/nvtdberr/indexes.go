package nvtdberr

import (
	"api/libs/nvt"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
)

// Querier is satisfied by *sql.DB, *sql.Conn, *sql.Tx and *goqu.Database.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// uniqueIndexQuery lists every UNIQUE and PRIMARY index of the current schema,
// one row per indexed column, in column order.
const uniqueIndexQuery = `SELECT TABLE_NAME, INDEX_NAME, COLUMN_NAME
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE() AND NON_UNIQUE = 0
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`

// indexes holds the current snapshot. nil means LoadIndexes has not succeeded yet.
var indexes atomic.Pointer[indexMap]

// LoadIndexes reads the unique-index metadata of the current MySQL/MariaDB
// schema so that duplicate-entry errors (1062), which only name an index, can
// report the real column in Info.Column.
//
// Call it once at startup after the database is reachable, and again after
// running migrations that add or rename unique indexes. It is a no-op for
// PostgreSQL, whose errors already carry the column.
//
// On failure the previous snapshot (if any) is kept, a warning is logged and the
// error is returned. Identify keeps working: it falls back to inferring the
// column from the index name, which is less reliable.
func LoadIndexes(ctx context.Context, q Querier) error {
	if dialect, ok := activeDialect(); ok && dialect == nvt.Postgres {
		return nil
	}
	return loadIndexes(ctx, q)
}

// loadIndexes is LoadIndexes without the dialect check.
func loadIndexes(ctx context.Context, q Querier) error {
	list, err := queryUniqueIndexes(ctx, q)
	if err != nil {
		log.Printf("nvtdberr: could not read unique-index metadata from information_schema.STATISTICS: %v. "+
			"Duplicate-key errors will fall back to inferring the column from the index name, which is less reliable. "+
			"Please check the database configuration (host, schema name) and make sure the connected user is "+
			"allowed to read information_schema.", err)
		return fmt.Errorf("nvtdberr: load unique indexes: %w", err)
	}

	if len(list) == 0 {
		log.Printf("nvtdberr: information_schema.STATISTICS returned no unique indexes for the current schema. " +
			"Please check that the connection selects the intended database and that the connected user has " +
			"privileges on its tables (information_schema only shows objects the user can access).")
	}
	indexes.Store(buildIndexMap(list))
	return nil
}

// queryUniqueIndexes runs uniqueIndexQuery. Rows of functional (expression)
// indexes have a NULL column and are skipped.
func queryUniqueIndexes(ctx context.Context, q Querier) ([]indexRow, error) {
	rows, err := q.QueryContext(ctx, uniqueIndexQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("nvtdberr: failed to close rows: %v", err)
		}
	}()

	var list []indexRow
	for rows.Next() {
		var table, index string
		var column sql.NullString
		if err := rows.Scan(&table, &index, &column); err != nil {
			return nil, err
		}
		if column.Valid {
			list = append(list, indexRow{table, index, column.String})
		}
	}
	return list, rows.Err()
}
