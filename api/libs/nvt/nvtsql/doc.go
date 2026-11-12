// Package nvtsql gives Norvet applications a configured *sql.DB without
// making them choose or wire a driver.
//
// Connect from the environment, in one call:
//
//	import (
//		"api/libs/nvt"
//		"api/libs/nvt/nvtsql"
//		_ "api/libs/nvt/nvtsql/drivers"
//	)
//
//	nvt.Init(nvt.Config{})
//
//	db, err := nvtsql.OpenFromEnv()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
// OpenFromEnv reads DB_DRIVER, the connection settings and the pool settings
// (see LoadConfig), opens the database with the matching driver, applies the
// pool and pings the server. To build the configuration yourself, fill a
// Config and call Open.
//
// # Drivers
//
// Drivers live in sub-packages that register themselves, so a program only
// links the databases it imports:
//
//	nvtsql/drivers    MySQL, MariaDB, PostgreSQL and SQLite
//	nvtsql/mysql      MySQL and MariaDB
//	nvtsql/postgres   PostgreSQL
//	nvtsql/sqlite     SQLite; also offers sqlite.Open(path) with no other dependency
//
// Open fails with ErrDriverNotLinked when the dialect's package was not
// imported, and with ErrUnsupportedDialect when DB_DRIVER is unset or unknown.
//
// # Adding a database
//
// Create a sub-package whose init registers an Opener for a new nvt.Dialect:
//
//	func init() { nvtsql.Register(nvt.SomeDialect, open) }
//
// # Layout
//
//	config.go   Config, LoadConfig and the environment variables
//	open.go     Register, Open, OpenFromEnv, pool and ping
package nvtsql
