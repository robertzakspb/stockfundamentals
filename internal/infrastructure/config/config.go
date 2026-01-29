// THIS IS AN ENV FILE FOR LOCAL DEVELOPMENT
package config

import (
	"os"
	"github.com/joho/godotenv"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

type Config struct {
	DB YDBConfig
}

type YDBConfig struct {
	ConnectionString string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("dev.env")
	if err != nil {
		logger.Log("FAILED LOAD THE ENVIRONMENT FILE " + "dev.env", logger.ALERT)
		return nil, err
	}

	dbConnectionString, ok := os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		panic("Unable to get the DB connection string")
	}
	ydbConfig := YDBConfig{
		dbConnectionString,
	}

	config := &Config{
		DB: ydbConfig,
	}

	return config, nil
}
