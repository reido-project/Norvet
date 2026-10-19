package config

type MySQLConfig struct {
	Charset string
	Loc     string
	TLS     string
}

func LoadMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		Charset: envString("MARIADB_CHARSET", "utf8mb4"),
		Loc:     envString("MARIADB_LOC", "Local"),
		TLS:     envString("MARIADB_TLS", "false"),
	}
}
