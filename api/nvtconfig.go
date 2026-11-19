package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/Norvet/nvtime"
	"gitlab.com/Norvet/nvtlog"
	"gitlab.com/Norvet/nvtsql"
)

func initNorvet() (diagnostic bool, err error) {
	// Flags
	diagnosticPtr := flag.Bool("D", false, "Run application in diagnostic mode")
	timePtr := flag.Bool("disable-time", true, "Disable time on terminal log")
	dirPtr := flag.String("dir", "logs/api", "Directory path for log output")
	flag.Parse()

	// nvtime Initialization
	if err := nvtime.Init(nvtime.Config{
		AppTimezone: os.Getenv("APP_TIMEZONE"),
		DBTimezone:  os.Getenv("DB_TIMEZONE"),
	}); err != nil {
		return false, err
	}

	// nvtlog Initialization
	if err := nvtlog.Init(nvtlog.Config{
		Color:      true,
		ShowTime:   *timePtr,
		Diagnostic: *diagnosticPtr,
		LogDir:     *dirPtr,
	}); err != nil {
		return false, err
	}

	database, err := databaseFromEnv()
	if err != nil {
		return false, err
	}

	// nvtsql Initialization
	if err := nvtsql.Init(map[uint]nvtsql.Database{nvtsql.DefaultID: database}); err != nil {
		return false, err
	}

	return *diagnosticPtr, nil
}

func databaseFromEnv() (nvtsql.Database, error) {
	var errs []error
	number := func(key string) int {
		raw := strings.TrimSpace(os.Getenv(key))
		if raw == "" {
			return 0
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s must be a whole number, got %q", key, raw))
		}
		return n
	}

	port := number("DB_PORT")
	pool := nvtsql.Pool{
		MaxOpenConns:    number("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    number("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: time.Duration(number("DB_CONN_MAX_LIFETIME")) * time.Second,
		ConnMaxIdleTime: time.Duration(number("DB_CONN_MAX_IDLE_TIME")) * time.Second,
	}
	if port < 0 || port > 65535 {
		errs = append(errs, fmt.Errorf("DB_PORT must be between 0 and 65535, got %d", port))
	}
	if len(errs) > 0 {
		return nvtsql.Database{}, errors.Join(errs...)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	var driver nvtsql.SQL
	switch strings.ToLower(strings.TrimSpace(os.Getenv("DB_DRIVER"))) {
	case "mariadb":
		driver = nvtsql.MariaDB{
			Host: host, Port: uint16(port), User: user, Password: password, Name: name,
			Charset: os.Getenv("MARIADB_CHARSET"), TLS: os.Getenv("MARIADB_TLS"),
		}
	case "mysql":
		driver = nvtsql.MySQL{
			Host: host, Port: uint16(port), User: user, Password: password, Name: name,
			Charset: os.Getenv("MARIADB_CHARSET"), TLS: os.Getenv("MARIADB_TLS"),
		}
	case "postgres", "postgresql":
		driver = nvtsql.Postgres{
			Host: host, Port: uint16(port), User: user, Password: password, Name: name,
			SSLMode: os.Getenv("POSTGRES_SSLMODE"),
		}
	case "sqlite", "sqlite3":
		driver = nvtsql.SQLite{Path: name}
	default:
		return nvtsql.Database{}, fmt.Errorf("DB_DRIVER must be mariadb, mysql, postgres or sqlite, got %q", os.Getenv("DB_DRIVER"))
	}

	return nvtsql.Database{Driver: driver, Pool: pool}, nil
}
