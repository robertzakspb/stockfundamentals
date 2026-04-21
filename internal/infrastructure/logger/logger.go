package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
)

var CURRENT_LOGGING_LEVELS = [4]LOG_LEVEL{INFORMATION, ERROR, ALERT, ERROR}

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

var mutex = sync.Mutex{}
var logFile *os.File

func writeLogToFile(config config.Config, message string, level LOG_LEVEL, logTime time.Time) error {
	// open input file
	mutex.Lock()
	defer mutex.Unlock()
	if logFile == nil {
		logFile, _ = os.OpenFile(config.Logger.FileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	log := []byte(logTime.String() + "." + level_name[level] + ": " + message + "\n")
	if _, err := logFile.Write(log); err != nil {
		fmt.Println(err)
	}

	return nil
}
