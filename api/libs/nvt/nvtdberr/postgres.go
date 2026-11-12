package nvtdberr

import (
	"api/libs/nvt/nvterr"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// postgresErrors maps PostgreSQL SQLSTATE codes (e.g. "23505") to client-safe nvterr errors.
// TODO: populate once the PostgreSQL error data is mapped.
var postgresErrors = map[string]*nvterr.Error{}

// identifyPostgres translates a PostgreSQL SQLSTATE into an nvterr error.
// Unmapped states return nvterr.ErrInternal.
func identifyPostgres(pe *pgconn.PgError) (*nvterr.Error, Info) {
	e, ok := postgresErrors[pe.Code]
	if !ok {
		e = nvterr.ErrInternal
	}
	info := Info{
		Database:   "", // not provided by the server
		Table:      pe.TableName,
		Column:     pe.ColumnName,
		Constraint: pe.ConstraintName,
	}
	// 23505 unique_violation: ColumnName is empty, but Detail is
	// "Key (username)=(myusername) already exists." -> take the part in the first parens.
	if pe.Code == "23505" && info.Column == "" {
		if i := strings.IndexByte(pe.Detail, '('); i >= 0 {
			if j := strings.IndexByte(pe.Detail[i:], ')'); j > 0 {
				info.Column = pe.Detail[i+1 : i+j] // may be "a, b" for composite keys
			}
		}
	}
	return e, info
}
