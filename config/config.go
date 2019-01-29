package config

import (
	"github.com/jackc/pgx"
)

type Config struct {
	Port string
	SqlConfig pgx.ConnPoolConfig
}

func GetConfig() (config Config) {
	config.Port = ":8080"
	config.SqlConfig = getSqlConfig()
	return
}

func getSqlConfig() (sqlConfig pgx.ConnPoolConfig) {
	var connConfig pgx.ConnConfig
	connConfig.Host = "localhost"
	connConfig.Port = 5432
	connConfig.User = "postgres"
	connConfig.Password = "postgres"
	connConfig.Database = "test"
	sqlConfig.ConnConfig = connConfig
	sqlConfig.MaxConnections = 6
	return
}