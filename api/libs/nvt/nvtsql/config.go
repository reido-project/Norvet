package nvtsql

import (
	"os"
	"strconv"
	"time"

	"api/libs/nvt"
)

// Config describes one database connection.
type Config struct {
	Dialect nvt.Dialect

	// Host, Port, User and Password locate and authenticate the server.
	// They are unused by SQLite.
	Host     string
	Port     string
	User     string
	Password string

	// Name is the database name, or the file path for SQLite.
	Name string

	// Location is the timezone for database timestamps. Nil means the
	// driver's default.
	Location *time.Location

	Pool Pool

	MySQL    MySQLOptions    // used by MySQL and MariaDB
	Postgres PostgresOptions // used by PostgreSQL
}

// Pool configures the connection pool. A zero field keeps the database/sql
// default.
type Pool struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// MySQLOptions are the MySQL and MariaDB specific settings.
type MySQLOptions struct {
	Charset string // for example utf8mb4
	TLS     string // true, false, skip-verify, preferred, or a registered name
}

// PostgresOptions are the PostgreSQL specific settings.
type PostgresOptions struct {
	SSLMode string // disable, require, verify-ca, verify-full, ...
}

// LoadConfig builds a Config from the Norvet context and the environment. It
// returns nvt.ErrNotInitialized if nvt.Init has not been called.
//
// The dialect and timezone come from the nvt context (DB_DRIVER and
// DB_TIMEZONE). The rest is read here:
//
//	DB_HOST                 default 127.0.0.1
//	DB_PORT                 default 3306 (MySQL, MariaDB) or 5432 (PostgreSQL)
//	DB_USER                 default root (MySQL, MariaDB) or postgres (PostgreSQL)
//	DB_PASSWORD
//	DB_NAME                 database name, or the file path for SQLite
//	DB_MAX_OPEN_CONNS       default 25
//	DB_MAX_IDLE_CONNS       default 10
//	DB_CONN_MAX_LIFETIME    seconds, default 1800
//	DB_CONN_MAX_IDLE_TIME   seconds, default 300
//	MARIADB_CHARSET         default utf8mb4 (also applies to MySQL)
//	MARIADB_TLS             default false (also applies to MySQL)
//	POSTGRES_SSLMODE        default disable
//
// Pool settings do not apply to SQLite, which always uses a single connection.
func LoadConfig() (Config, error) {
	ctx := nvt.Ctx()
	if ctx == nil {
		return Config{}, nvt.ErrNotInitialized
	}

	cfg := Config{
		Dialect:  ctx.Database,
		Host:     envString("DB_HOST", "127.0.0.1"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Location: ctx.DBTimezone,
		MySQL: MySQLOptions{
			Charset: envString("MARIADB_CHARSET", "utf8mb4"),
			TLS:     envString("MARIADB_TLS", "false"),
		},
		Postgres: PostgresOptions{
			SSLMode: envString("POSTGRES_SSLMODE", "disable"),
		},
	}

	switch ctx.Database {
	case nvt.Postgres:
		cfg.Port = envString("DB_PORT", "5432")
		cfg.User = envString("DB_USER", "postgres")
	default:
		cfg.Port = envString("DB_PORT", "3306")
		cfg.User = envString("DB_USER", "root")
	}

	if ctx.Database != nvt.SQLite {
		cfg.Pool = Pool{
			MaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(envInt("DB_CONN_MAX_LIFETIME", 1800)) * time.Second,
			ConnMaxIdleTime: time.Duration(envInt("DB_CONN_MAX_IDLE_TIME", 300)) * time.Second,
		}
	}

	return cfg, nil
}

func envString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// envInt reads a positive integer; anything else gives fallback.
func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
