package aws

import (
  "errors"
  "strings"
  "github.com/PyramidSystemsInc/go/commands"
)

func GetEcrUrl() (string, error) {
  output := commands.Run("aws ecr get-login", "")
  words := strings.Split(output, " ")
  url := words[len(words) - 1]
  if strings.HasPrefix(url, "https://") {
    return url, nil
  } else {
    return "", errors.New("URL not found. Are your AWS credentials configured?")
  }
}
