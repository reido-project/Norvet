# nvtdberr

Translates database errors into client-safe [`nvterr`](../nvterr) errors.

A raw driver error such as `Error 1062: Duplicate entry 'bob' for key 'users.username'`
must never reach a client: it leaks schema and row values. `nvtdberr.Identify`
turns it into an `*nvterr.Error` with a stable code and HTTP status, keeps the
original as the hidden cause (visible in logs only), and tells you which table,
column or constraint was involved.

```go
import "api/libs/nvt/nvtdberr"
```

## Usage

```go
u, err := s.repo.Create(ctx, user)
if err != nil {
	e, info := nvtdberr.Identify(err)

	// Specific, client-friendly errors for the cases you care about:
	if errors.Is(e, nvterr.ErrDuplicateValue) {
		switch info.Column {
		case "email":
			return nil, nvterr.ErrEmailAlreadyRegistered.WithCause(err)
		case "username":
			return nil, nvterr.ErrUsernameAlreadyRegistered.WithCause(err)
		}
	}
	return nil, e // everything else: a safe generic error
}
```

`Identify` never returns nil for a non-nil error, so `return nil, e` is always
safe.

### Startup

```go
loadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
_ = nvtdberr.LoadIndexes(loadCtx, db) // failures are logged; Identify keeps working
cancel()
```

MySQL and MariaDB report a duplicate key by **index name**, not column.
`LoadIndexes` reads the unique-index metadata so `Info.Column` is exact (also
for composite keys: `"org_id, email"`). Call it again after migrations that add
or rename unique indexes. Without it the column is inferred from the index name,
which is wrong for names that do not follow a convention, and a warning is
logged once per index. It is a no-op on PostgreSQL.

## What you get back

| Input | Result |
|---|---|
| `nil` | `nil`, `Info{}` |
| an `*nvterr.Error` in the chain, `context.Canceled`, `context.DeadlineExceeded` | handled by `nvterr.From` |
| MariaDB / MySQL server error (`*mysql.MySQLError`) | the dialect's table, e.g. 1062 gives `duplicate_value` (422), 1146 gives `not_found` (404) |
| PostgreSQL server error (`*pgconn.PgError`) | the PostgreSQL table, and `Info` from the error fields |
| `sql.ErrNoRows` | `not_found` (404). Prefer `nvterr.NewNotFound("invoice")` at the call site for a better message |
| bad / closed connection | `database_unavailable` (503) |
| network timeout | `database_timeout` (504) |
| unmapped server error number | `internal` (500) |
| anything else | `internal` (500), and a one-time warning that it is not a database error |

The original error is always reachable with `errors.Is`/`errors.As` and appears
in `e.Error()`, never in `e.Message()` or the JSON.

### `Info`

```go
type Info struct {
	Database, Table, Column, Constraint string
}
info.Subject() // most specific non-empty of Column, Table, Database
```

Fields are empty when the database error does not say. `Info` never contains row
values.

## Dialects

The dialect for MariaDB and MySQL comes from `DB_DRIVER` through the nvt
context. PostgreSQL errors are recognized by type and need no configuration.

| Dialect | Error table | Status |
|---|---|---|
| MariaDB | `mariadb.go`, about 620 server and client error numbers | complete |
| MySQL | `mysql.go` | **in progress**: all MySQL server errors map to `internal` for now |
| PostgreSQL | `postgres.go` | **in progress**: all PostgreSQL server errors map to `internal` for now (`Info` already works) |

If a MySQL/MariaDB error arrives while the dialect is unset or invalid, it maps
to `internal` and a warning naming `DB_DRIVER` is logged once.

## Maintaining

- Add or change entries in the dialect's table (`mariadbErrors`, ...). Every
  entry must be a non-500 error; unlisted numbers already fall back to
  `nvterr.ErrInternal`. A test enforces this.
- To teach `Info` about another error number, add a parser to `mysqlSubjects` in
  `mysql_subject.go`.
- File layout is described in [doc.go](doc.go).
