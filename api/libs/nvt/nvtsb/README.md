# sqlbuilder

Connects to the database from the environment and returns a query builder. Currently supports [goqu](https://github.com/doug-martin/goqu).

```go
nvt.Init(nvt.Config{})

db, err := sqlbuilder.ConnectGoqu()
if err != nil {
    log.Fatal(err)
}
defer sqlbuilder.CloseGoqu(db)

var users []User
err = db.From("users").Where(goqu.C("active").IsTrue()).ScanStructs(&users)
```

| Function | Purpose |
|---|---|
| `ConnectGoqu()` | open the database from the environment and wrap it |
| `WrapGoqu(db, dialect)` | wrap a `*sql.DB` you already have (for example from `nvtsql.Open`) |
| `CloseGoqu(db)` | close the underlying connection pool |

Connection settings, pool settings and supported databases (MySQL, MariaDB, PostgreSQL, SQLite) are those of [`nvtsql`](../nvtsql), whose README lists every environment variable. `nvt.Init` must run first; otherwise `ConnectGoqu` returns `nvt.ErrNotInitialized`.

## Adding a builder

Add a file with a `Wrap<Builder>` function and a one-line `Connect<Builder>`. Nothing existing changes:

```go
func ConnectSquirrel() (squirrel.StatementBuilderType, error) {
    return connect(WrapSquirrel)
}
```

`connect` opens the database, calls your wrap function and closes the database again if wrapping fails.
