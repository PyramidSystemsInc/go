package commands

import (
  "os/exec"
  "strings"
  "github.com/PyramidSystemsInc/go/directories"
  "github.com/PyramidSystemsInc/go/errors"
)

// Runs a command as if ran from the terminal
func Run(fullCommand string, directory string) {
  command, arguments := separateCommand(fullCommand)
  cmd := exec.Command(command, arguments...)
  if directory == "" {
    cmd.Dir = directories.GetWorking()
  } else {
    cmd.Dir = directory
  }
  err := cmd.Run()
  errors.LogIfError(err)
}

func separateCommand(fullCommand string) (string, []string) {
  split := strings.Split(fullCommand, " ")
  return split[0], split[1:len(split)]
}
