// THIS IS AN ENV FILE FOR LOCAL DEVELOPMENT
package config

import (
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DB YDBConfig
}

type YDBConfig struct {
	ConnectionString string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err

	}


	dbConnectionString, ok := os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		panic("Unable to get the DB connection string")
	}
	ydbConfig := YDBConfig {
		dbConnectionString,
	}

	config := &Config{
		DB: ydbConfig,
	}

	return config, nil
}