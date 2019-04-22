package errors

import (
  "errors"
  "os"

  "github.com/PyramidSystemsInc/go/logger"
)

// LogAndQuit - Log a custom string as an error and halt execution
func LogAndQuit(message string) {
  logger.Err(message)
  os.Exit(-1)
}

// LogIfError - Logs the error and continues execution
func LogIfError(err error) {
  if err != nil {
    logger.Err(err.Error())
  }
}

// QuitIfError - Logs the error and halts execution of the program
func QuitIfError(err error) {
  if err != nil {
    logger.Err(err.Error())
    os.Exit(-1)
  }
}

// ReturnIfError - Logs the error and returns out of the current function
func ReturnIfError(err error) {
  if err != nil {
    logger.Err(err.Error())
    return
  }
}

// New - Returns a new error
func New(err string) error {
  return errors.New(err)
}
