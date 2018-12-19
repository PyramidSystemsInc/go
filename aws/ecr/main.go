package ecr

import (
  "errors"
  "strings"
  "github.com/PyramidSystemsInc/go/commands"
)

func GetUrl() (string, error) {
  output := commands.Run("aws ecr get-login", "")
  words := strings.Split(output, " ")
  url := words[len(words) - 1]
  if strings.HasPrefix(url, "https://") {
    return strings.TrimLeft(url, "https://"), nil
  } else {
    return "", errors.New("URL not found. Are your AWS credentials configured?")
  }
}

func Login(region string) {
  output := commands.Run("aws ecr get-login --no-include-email --region " + region, "")
  commands.Run(output, "")
}
