# nvtsql

A configured `*sql.DB` for MySQL, MariaDB, PostgreSQL and SQLite, driven by environment variables. No query builder, no driver wiring.

```go
import (
    "api/libs/nvt"
    "api/libs/nvt/nvtsql"
    _ "api/libs/nvt/nvtsql/drivers" // or only the ones you need
)

nvt.Init(nvt.Config{})

db, err := nvtsql.OpenFromEnv()
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

Need a query builder? Use [`sqlbuilder`](../sqlbuilder), which sits on top of this package. Need only a SQLite file and nothing else?

```go
import "api/libs/nvt/nvtsql/sqlite"

db, err := sqlite.Open("logs/app.sqlite") // database/sql + the SQLite driver only
```

## Environment

`DB_DRIVER` (read by `nvt.Init`) selects the database: `mariadb`, `mysql`, `postgres` (or `postgresql`), `sqlite` (or `sqlite3`). `DB_TIMEZONE` (also read by `nvt.Init`) sets the timezone of database timestamps.

| Variable | Default | Notes |
|---|---|---|
| `DB_HOST` | `127.0.0.1` | |
| `DB_PORT` | `3306` / `5432` | `5432` for PostgreSQL |
| `DB_USER` | `root` / `postgres` | `postgres` for PostgreSQL |
| `DB_PASSWORD` | | any characters; no escaping needed |
| `DB_NAME` | | database name; the file path for SQLite |
| `DB_MAX_OPEN_CONNS` | `25` | not used by SQLite (always 1) |
| `DB_MAX_IDLE_CONNS` | `10` | |
| `DB_CONN_MAX_LIFETIME` | `1800` | seconds |
| `DB_CONN_MAX_IDLE_TIME` | `300` | seconds |
| `MARIADB_CHARSET` | `utf8mb4` | MySQL and MariaDB |
| `MARIADB_TLS` | `false` | `true`, `false`, `skip-verify`, `preferred` |
| `POSTGRES_SSLMODE` | `disable` | |

Invalid or non-positive numbers fall back to the default.

## Choosing drivers

| Import | Links |
|---|---|
| `nvtsql/drivers` | all four databases |
| `nvtsql/mysql` | MySQL, MariaDB |
| `nvtsql/postgres` | PostgreSQL |
| `nvtsql/sqlite` | SQLite (needs cgo) |

## Errors

`Open` returns a wrapped `nvtsql.ErrUnsupportedDialect` when `DB_DRIVER` is unset or unknown, `nvtsql.ErrDriverNotLinked` when the matching package was not imported, and `nvt.ErrNotInitialized` (from `LoadConfig`) when `nvt.Init` was not called. Check them with `errors.Is`. A server that does not answer within 10 seconds fails the ping.

## Adding a database

Create a sub-package whose `init` calls `nvtsql.Register(dialect, opener)`; see `nvtsql/postgres` for a complete example.
