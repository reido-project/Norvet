package config

import (
	"os"
	"strconv"
)

type Config struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int

	MySQL    *MySQLConfig
	Postgres *PostgresConfig
}

func LoadConfig() Config {
	driver := envString("DB_DRIVER", "mariadb")

	cfg := Config{
		Driver:          driver,
		Host:            envString("DB_HOST", "127.0.0.1"),
		Port:            envString("DB_PORT", "3306"),
		User:            envString("DB_USER", "root"),
		Password:        os.Getenv("DB_PASSWORD"),
		Name:            os.Getenv("DB_NAME"),
		MaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: envInt("DB_CONN_MAX_LIFETIME", 1800),
		ConnMaxIdleTime: envInt("DB_CONN_MAX_IDLE_TIME", 300),
	}

	switch driver {
	case "mysql", "mariadb":
		cfg.MySQL = LoadMySQLConfig()
	case "postgres", "postgresql":
		cfg.Postgres = LoadPostgresConfig()
	}

	return cfg
}

func envString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
