package commands

import (
  "io/ioutil"
  "os/exec"
  "strings"
  "regexp"
  "github.com/PyramidSystemsInc/go/directories"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/logger"
  "github.com/PyramidSystemsInc/go/str"
)

// Run - Runs a command as if ran from the terminal
func Run(fullCommand string, directory string) (string, error) {
  command, arguments := separateCommand(fullCommand)
  cmd := exec.Command(command, arguments...)

  stdout, err := cmd.StdoutPipe()
  errors.LogIfError(err)
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
  errOutput, err := ioutil.ReadAll(stderr)
  errors.LogIfError(err)
  err = cmd.Wait()
  if err != nil {
    err = errors.New(str.Concat(err.Error(), ": ", strings.TrimRight(string(errOutput), "\n")))
  }
  out := strings.TrimRight(string(output), "\n")
  return out, err
}

// RunWithStdin - Runs a command as if ran from the terminal
func RunWithStdin(fullCommand string, data string, directory string) string {
  command, arguments := separateCommand(fullCommand)
  cmd := exec.Command(command, arguments...)
  cmd.Stdin = strings.NewReader(data)
  stdout, err := cmd.StdoutPipe()
  errors.LogIfError(err)
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


func replaceRelativeWithFullPath(directory string) string {
  return strings.Replace(directory, ".", directories.GetWorking(), 1)
}

/*
 * Separates a Command (executable) from its arguments. Delimitters taken into account are spaces
 * but spaces within double quotes are conserved. Double quotes are eventually removed as well since
 * GOlang escapes them prior to sending to the command line.
 * @return string - command
 * @return []string - array of arguments
 */
func separateCommand(text string) (string, []string) {
  //Matches words seperated by a space except if enclosed in double quotes
  exp := regexp.MustCompile(`[^\s"']+|"[^"]*"|'[^']`)
  expResults := exp.FindAllString(text, -1)

  sanitizedVals := []string{}

  for counter := range expResults {
    currSanitizedValue := strings.Replace(expResults[counter],`"`, "", -1)
    sanitizedVals = append(sanitizedVals, currSanitizedValue)
  }
  return sanitizedVals[0], sanitizedVals[1:len(sanitizedVals)]
}

