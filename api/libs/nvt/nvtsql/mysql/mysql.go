// Package mysql registers the MySQL and MariaDB connector with nvtsql.
// Import it for its side effect:
//
//	import _ "api/libs/nvt/nvtsql/mysql"
package mysql

import (
	"database/sql"
	"fmt"
	"net"

	driver "github.com/go-sql-driver/mysql"

	"api/libs/nvt"
	"api/libs/nvt/nvtsql"
)

func init() {
	nvtsql.Register(nvt.MySQL, open)
	nvtsql.Register(nvt.MariaDB, open)
}

func open(cfg nvtsql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("nvtsql: open %s: %w", cfg.Dialect, err)
	}
	return db, nil
}

// dsn builds the DSN with the driver's own formatter. Credentials must not be
// URL-escaped: go-sql-driver/mysql takes them literally.
func dsn(cfg nvtsql.Config) string {
	c := driver.NewConfig()
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.Net = "tcp"
	c.Addr = net.JoinHostPort(cfg.Host, cfg.Port)
	c.DBName = cfg.Name
	c.ParseTime = true
	c.TLSConfig = cfg.MySQL.TLS
	if cfg.Location != nil {
		c.Loc = cfg.Location
	}
	if cfg.MySQL.Charset != "" {
		c.Params = map[string]string{"charset": cfg.MySQL.Charset}
	}
	return c.FormatDSN()
}
