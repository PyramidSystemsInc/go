package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	ERR  LogLevel = 0
	WARN LogLevel = 1
	INFO LogLevel = 2
)

var Levels = map[LogLevel]string{
    ERR:  "ERROR",
    WARN: "WARNING",
    INFO: "INFO",
}

var logLevel = ERR

// Err - Logs and error from a string
func Err(message ...string) {
	// Always log errors
	log(ERR, message...)
}

// Info - Logs output as information
func Info(message ...string) {
	if logLevel >= INFO {
		log(INFO, message...)
	}
}

// SetLogLevel - Sets the maximum level of logging desired. All types more frequent than the type specified will not be output
func SetLogLevel(level LogLevel) {
	logLevel = level
}

// Warn - Logs output as a warning
func Warn(message ...string) {
	if logLevel >= WARN {
		log(WARN, message...)
	}
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func log(typ LogLevel, message ...string) {
	if typ == ERR {
		fmt.Fprintf(os.Stderr, "[ %s]  %s: %s\n", timestamp(), Levels[typ], strings.Join(message, " "))
	} else {
		fmt.Fprintf(os.Stdout, "[ %s]  %s: %s\n", timestamp(), Levels[typ], strings.Join(message, " "))
	}
}

func ParseLevel(userInput string) (LogLevel, bool) {
	var defaultLog LogLevel
	for k, v := range Levels {
		if userInput == v {
			return k, true
		}
	}
	return defaultLog, false
}
