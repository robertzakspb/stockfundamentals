package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     YDBConfig
	Logger LoggerConfig
}

type LoggerConfig struct {
	FileLocation string //Location of the log file in the file system
	Mode         LogMode
}

type LogMode string

const (
	CONSOLE LogMode = "CONSOLE"
	FILE    LogMode = "FILE"
)

var logModeFromString = map[string]LogMode{
	"CONSOLE": CONSOLE,
	"FILE":    FILE,
}

type YDBConfig struct {
	ConnectionString string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("dev.env")
	if err != nil {
		panic("ALERT: FAILED LOAD THE ENVIRONMENT FILE " + "dev.env")
	}

	dbConnectionString, found := os.LookupEnv("DB_CONNECTION_STRING")
	if !found {
		panic("Unable to get the DB connection string")
	}

	ydbConfig := YDBConfig{
		dbConnectionString,
	}

	logMode, found := os.LookupEnv("LOG_MODE")
	if !found {
		logMode = "CONSOLE"
		fmt.Println("Unexpectedly did not find the logging mode in the .env file")
	}

	fileLocation, found := os.LookupEnv("FILE_LOCATION")
	if !found && logMode == string(FILE) {
		panic("Unable to get the the log file location despite the log mode being set to FILE")
	}

	loggerConfig := LoggerConfig{
		Mode:         logModeFromString[logMode],
		FileLocation: fileLocation,
	}

	config := &Config{
		DB:     ydbConfig,
		Logger: loggerConfig,
	}

	return config, nil
}
