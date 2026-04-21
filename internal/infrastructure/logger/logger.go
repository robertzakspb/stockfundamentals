package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
)

var CURRENT_LOGGING_LEVELS = [4]LOG_LEVEL{ERROR, ALERT, ERROR}

type LOG_LEVEL int

const (
	INFORMATION LOG_LEVEL = iota
	WARNING
	ERROR
	ALERT
)

var level_name = map[LOG_LEVEL]string{
	INFORMATION: "INFORMATION",
	WARNING:     "WARNING",
	ERROR:       "ERROR",
	ALERT:       "ALERT",
}

func Log(message string, level LOG_LEVEL) {
	logTime := time.Now()
	shouldLog := false
	for _, l := range CURRENT_LOGGING_LEVELS {
		if level == l {
			shouldLog = true
		}
	}

	if !shouldLog {
		return
	}

	configuration, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	if configuration.Logger.Mode == config.CONSOLE {
		fmt.Println(message)
	}

	if configuration.Logger.Mode == config.FILE {
		go writeLogToFile(*configuration, message, level, logTime)
	}
}

func writeLogToFile(config config.Config, message string, level LOG_LEVEL, logTime time.Time) error {
	// open input file
	logFile, err := os.Open(config.Logger.FileLocation)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Println("ALERT: Failed to close the log fie")
		}
	}()

	log := []byte(logTime.String() + level_name[level] + message)
	if _, err := logFile.Write(log); err != nil {
		fmt.Println(err)
	}

	return nil
}
