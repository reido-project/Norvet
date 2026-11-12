// Package nvtdberr translates database errors into client-safe nvterr errors.
//
// A raw driver error ("Error 1062: Duplicate entry 'bob' for key
// 'users.username'") is never something to show a client. Identify turns it
// into an *nvterr.Error with a stable code and status, keeps the original as
// the hidden cause for logs, and reports which table, column or constraint was
// involved so the caller can build a precise message:
//
//	user, err := repo.Create(ctx, u)
//	if err != nil {
//		e, info := nvtdberr.Identify(err)
//		if errors.Is(e, nvterr.ErrDuplicateValue) && info.Column == "email" {
//			return nil, nvterr.ErrEmailAlreadyRegistered.WithCause(err)
//		}
//		return nil, e
//	}
//
// # Dialects
//
// MariaDB, MySQL and PostgreSQL are supported. Each has its own lookup table
// (mariadb.go, mysql.go, postgres.go). PostgreSQL errors are recognized by
// their type. MySQL and MariaDB share one error type, so the dialect comes
// from the active nvt context (DB_DRIVER). The MySQL and PostgreSQL tables are
// still being populated: until then their server errors map to
// nvterr.ErrInternal.
//
// # Duplicate keys
//
// MySQL and MariaDB report a duplicate key by index name, not by column. Call
// LoadIndexes once at startup (and again after migrations that change unique
// indexes) so Info.Column is exact; without it the column is inferred from the
// index name, which is less reliable and logged as a warning once per index.
//
// # Layout
//
//	identify.go        Identify and its resolution order
//	driver.go          database/sql and connection failures without a server code
//	dialect.go         the only place that reads the nvt context
//	info.go            Info
//	mariadb.go         MariaDB error-number table
//	mysql.go           MySQL error-number table, lookup shared with MariaDB
//	postgres.go        PostgreSQL SQLSTATE table
//	mysql_subject.go   extracts table/column/constraint from MySQL/MariaDB messages
//	indexes.go         LoadIndexes and the current snapshot
//	index_map.go       the immutable (table, index) to columns map
//	column_resolve.go  index to column resolution, with inference fallback
//	warn.go            warn-once logging
package nvtdberr
