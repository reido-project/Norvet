package nvtdberr

import (
	"api/libs/nvt"
	"api/libs/nvt/nvterr"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

// Identify converts err, returned by a database call, into an *nvterr.Error.
//
// Resolution order:
//
//  1. nil returns nil.
//  2. An *nvterr.Error already in the chain, context.Canceled and
//     context.DeadlineExceeded are handled by nvterr.From.
//  3. A server error carrying an error code is looked up in the table of its
//     dialect: PostgreSQL errors by their type, MySQL and MariaDB errors by the
//     dialect of the active nvt context. Unmapped codes return
//     nvterr.ErrInternal.
//  4. Driver and connection failures with no server code (bad connection,
//     timeouts, sql.ErrNoRows, ...) are mapped by identifyDriver.
//  5. Anything else is not a database error. A warning is logged once per
//     error type for the developer and nvterr.ErrInternal is returned so the
//     runtime keeps going.
//
// The returned Info names the database objects the error is about (table,
// column, constraint, ...) when the dialect's error exposes them; fields are
// empty otherwise. It never contains row values.
//
// The original error is attached with WithCause, so it appears in logs via
// Error() but is never exposed to clients.
func Identify(err error) (*nvterr.Error, Info) {
	dialect, _ := activeDialect()
	return identify(dialect, err)
}

// identify is Identify with the dialect passed in, so it can be exercised
// without the global nvt context.
func identify(dialect nvt.Dialect, err error) (*nvterr.Error, Info) {
	if err == nil {
		return nil, Info{}
	}
	if nvterr.IsKnown(err) {
		return nvterr.From(err), Info{}
	}

	if e, info, ok := identifyServer(dialect, err); ok {
		return e.WithCause(err), info
	}
	if e, ok := identifyDriver(err); ok {
		return e.WithCause(err), Info{}
	}

	warnOnce(fmt.Sprintf("not-a-db-error:%T", err),
		"nvtdberr: received an error that is not a database error (%T). "+
			"nvtdberr.Identify only accepts errors returned by database calls; "+
			"convert other errors with nvterr.From or return an nvterr error directly. Original error: %v", err, err)
	return nvterr.ErrInternal.WithCause(err), Info{}
}

// identifyServer maps errors reported by the database server itself. ok is
// false when err is not a server error.
func identifyServer(dialect nvt.Dialect, err error) (e *nvterr.Error, info Info, ok bool) {
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		e, info = identifyPostgres(pe)
		return e, info, true
	}

	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return nil, Info{}, false
	}
	switch dialect {
	case nvt.MariaDB:
		e, info = identifyMariaDB(me)
	case nvt.MySQL:
		e, info = identifyMySQL(me)
	default:
		warnOnce("dialect:"+string(dialect),
			"nvtdberr: got a MySQL/MariaDB error but the active database dialect is %q, so it cannot be mapped. "+
				"Check that nvt.Init was called and that the DB_DRIVER environment variable is mariadb or mysql; returning an internal error for: %v",
			dialect, err)
		e = nvterr.ErrInternal
	}
	return e, info, true
}
