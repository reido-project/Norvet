package config

type PostgresConfig struct {
	SSLMode  string
	Timezone string
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		SSLMode:  envString("POSTGRES_SSLMODE", "disable"),
		Timezone: envString("POSTGRES_TIMEZONE", "Asia/Jakarta"),
	}
}
