package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type LogLevel struct {
	Info bool
	Warn bool
	Err  bool
}

var logLevel LogLevel

// Err - Logs and error from a string
func Err(message ...string) {
	if logLevel.Err {
		fmt.Fprintln(os.Stderr, "[ "+timestamp()+"   ERROR ]: "+strings.Join(message, " "))
	}
}

// Info - Logs output as information
func Info(message ...string) {
	if logLevel.Info {
		fmt.Println("[ " + timestamp() + "    INFO ]: " + strings.Join(message, " "))
	}
}

// SetLogLevel - Sets the maximum level of logging desired. All types more frequent than the type specified will not be output
func SetLogLevel(level string) {
	if level == "info" || level == "warn" || level == "err" {
		logLevel = LogLevel{
			Info: level == "info",
			Warn: level == "info" || level == "warn",
			Err:  true,
		}
	} else {
		fmt.Println("The valid log types are: 'info', 'warn', and 'err'")
	}
}

// Warn - Logs output as a warning
func Warn(message ...string) {
	if logLevel.Warn {
		fmt.Println("[ " + timestamp() + " WARNING ]: " + strings.Join(message, " "))
	}
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
