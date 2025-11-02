package logger

import (
	"fmt"
)

var CURRENT_LOGGING_LEVELS = [...]LOG_LEVEL{ERROR, INFORMATION, ALERT, ERROR}

type LOG_LEVEL int

const (
	INFORMATION LOG_LEVEL = iota
	WARNING
	ERROR
	ALERT
)

func Log(message string, level LOG_LEVEL) error {
	shouldLog := false
	for _, l := range CURRENT_LOGGING_LEVELS {
		if level == l {
			shouldLog = true
		}
	}

	if shouldLog {
		fmt.Println(message)
	}

	return nil
}