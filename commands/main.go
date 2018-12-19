package commands

import (
  "io/ioutil"
  "os/exec"
  "strings"
  "github.com/PyramidSystemsInc/go/directories"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/logger"
)

// Runs a command as if ran from the terminal
func Run(fullCommand string, directory string) string {
  command, arguments := separateCommand(fullCommand)
  cmd := exec.Command(command, arguments...)
  stdout, err := cmd.StdoutPipe()
  stderr, err := cmd.StderrPipe()
  errors.LogIfError(err)
  if directory == "" {
    cmd.Dir = directories.GetWorking()
  } else if strings.HasPrefix(directory, "./") {
    cmd.Dir = replaceRelativeWithFullPath(directory)
  } else {
    cmd.Dir = directory
  }
  err = cmd.Start()
  errors.LogIfError(err)
  output, err := ioutil.ReadAll(stdout)
  errors.LogIfError(err)
  errorOutput, err := ioutil.ReadAll(stderr)
  errors.LogIfError(err)
  if string(errorOutput) != "" {
    logger.Warn(string(errorOutput))
  }
  err = cmd.Wait()
  errors.LogIfError(err)
  out := strings.TrimRight(string(output), "\n")
  return out
}

func separateCommand(fullCommand string) (string, []string) {
  split := strings.Split(fullCommand, " ")
  return split[0], split[1:len(split)]
}

func replaceRelativeWithFullPath(directory string) string {
  return strings.Replace(directory, ".", directories.GetWorking(), 1)
}
