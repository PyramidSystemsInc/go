package directories

import (
  "os"
  "os/user"
  "github.com/PyramidSystemsInc/go/errors"
)

// Creates a directory (if it does not already exist)
func Create(directory string) {
  err := os.MkdirAll(directory, os.ModePerm)
  errors.QuitIfError(err)
}

// Returns the home directory of the user running the program
func GetHome() string {
  user, err := user.Current()
  errors.QuitIfError(err)
  return user.HomeDir
}

// Returns the working directory
func GetWorking() string {
  workingDirectory, err := os.Getwd()
  errors.LogIfError(err)
  return workingDirectory
}
