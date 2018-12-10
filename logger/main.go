package logger

import (
  "fmt"
  "time"
)

type LogLevel struct {
  Info  bool
  Warn  bool
  Err   bool
}

var logLevel LogLevel

func Err(message string) {
  if (logLevel.Err) {
    fmt.Println("[ " + timestamp() + "   ERROR ]: " + message)
  }
}

func Info(message string) {
  if (logLevel.Info) {
    fmt.Println("[ " + timestamp() + "    INFO ]: " + message)
  }
}

func SetLogLevel(level string) {
  if (level == "info" || level == "warn" || level == "err") {
    logLevel = LogLevel{
      Info:  level == "info",
      Warn:  level == "info" || level == "warn",
      Err:   true,
    }
  } else {
    fmt.Println("The valid log types are: 'info', 'warn', and 'err'")
  }
}

func Warn(message string) {
  if (logLevel.Warn) {
    fmt.Println("[ " + timestamp() + " WARNING ]: " + message)
  }
}

func timestamp() string {
  return time.Now().UTC().Format(time.RFC3339)
}
