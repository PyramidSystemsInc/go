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

var logLevel = ERR

func log(typ LogLevel, message ...string) {
	// msg := "[ " + timestamp() + "]  %s: %s"
	if typ == ERR {
		fmt.Fprintf(os.Stderr, "[ %s]  %s: %s\n", timestamp(), typ, strings.Join(message, " "))
	} else {
		fmt.Fprintf(os.Stdout, "[ %s]  %s: %s\n", timestamp(), typ, strings.Join(message, " "))
	}
}

// Err - Logs and error from a string
func Err(message ...string) {
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

func (lvl LogLevel) String() string {
	names := [...]string{"ERROR!", "WARNING", "INFO"}
	return names[lvl]
}
