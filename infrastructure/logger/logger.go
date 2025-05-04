package logger

import "fmt"

type LOG_LEVEL int
const (
	INFORMATION LOG_LEVEL = iota
	WARNING 
	ERROR 
	ALERT
)

func Log(message string, level LOG_LEVEL) error {
	fmt.Println(message)

	return nil
}