package commands

import (
  "os/exec"
  "strings"
  "github.com/PyramidSystemsInc/go/directories"
  "github.com/PyramidSystemsInc/go/errors"
)

// Runs a command as if ran from the terminal
func Run(fullCommand string, directory string) string {
  command, arguments := separateCommand(fullCommand)
  cmd := exec.Command(command, arguments...)
  if directory == "" {
    cmd.Dir = directories.GetWorking()
  } else if strings.HasPrefix(directory, "./") {
    cmd.Dir = replaceRelativeWithFullPath(directory)
  } else {
    cmd.Dir = directory
  }
  out, _ := cmd.Output()
  err := cmd.Run()
  harmlessError := "exec: already started"
  if err.Error() == harmlessError {
    return string(out)
  } else {
    errors.LogIfError(err)
    return string(out)
  }
}

func separateCommand(fullCommand string) (string, []string) {
  split := strings.Split(fullCommand, " ")
  return split[0], split[1:len(split)]
}

func replaceRelativeWithFullPath(directory string) string {
  return strings.Replace(directory, ".", directories.GetWorking(), 1)
}
