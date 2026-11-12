package nvt

import "strings"

// Dialect identifies a supported SQL database.
type Dialect string

const (
	MariaDB  Dialect = "mariadb"
	MySQL    Dialect = "mysql"
	Postgres Dialect = "postgres"
	SQLite   Dialect = "sqlite"
)

// ParseDialect converts a DB_DRIVER value to a Dialect. Matching ignores case
// and surrounding spaces and accepts the common aliases "postgresql" and
// "sqlite3". ok is false for anything else.
func ParseDialect(name string) (d Dialect, ok bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "mariadb":
		return MariaDB, true
	case "mysql":
		return MySQL, true
	case "postgres", "postgresql":
		return Postgres, true
	case "sqlite", "sqlite3":
		return SQLite, true
	}
	return "", false
}

// String returns the canonical lower-case name.
func (d Dialect) String() string { return string(d) }
