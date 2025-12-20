package config

import "os"

type EnvConfig struct {
	AppPort       string
	PostgresqlURL string
	RDBAddress    string
	RDBPassword   string
}

var LocalConfig *EnvConfig

func initialConfig() *EnvConfig {
	return &EnvConfig{
		AppPort:       os.Getenv("APP_PORT"),
		PostgresqlURL: os.Getenv("POSTGRESQL_URL"),
		RDBAddress: os.Getenv("RDB_ADDR"),
		RDBPassword: os.Getenv("RDB_PASS"),
	}
}

func SetConfig() {
	LocalConfig = initialConfig()
}
