package provider

import (
	"api/persistence/config"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

func OpenMySQL(cfg config.Config) (*sql.DB, string, error) {
	if cfg.MySQL == nil {
		return nil, "", fmt.Errorf("mysql config is missing")
	}

	encodedUser := url.QueryEscape(cfg.User)
	encodedPassword := url.QueryEscape(cfg.Password)
	encodedDBName := url.PathEscape(cfg.Name)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=%s&loc=%s&tls=%s",
		encodedUser,
		encodedPassword,
		cfg.Host,
		cfg.Port,
		encodedDBName,
		url.QueryEscape(cfg.MySQL.Charset),
		url.QueryEscape(cfg.MySQL.Loc),
		url.QueryEscape(cfg.MySQL.TLS),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("open MySQL: %w", err)
	}

	return db, "mysql", nil
}
