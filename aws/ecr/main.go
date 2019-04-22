package ecr

import (
  "strings"
  "github.com/PyramidSystemsInc/go/commands"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

// GetUrl - Returns the URL of your ECR repository
func GetUrl() (string, error) {
  output, err := commands.Run("aws ecr get-login", "")
  errors.LogIfError(err)
  words := strings.Split(output, " ")
  url := words[len(words) - 1]
  if strings.HasPrefix(url, "https://") {
    return strings.TrimLeft(url, "https://"), nil
  } else {
    return "", errors.New("URL not found. Are your AWS credentials configured?")
  }
}

// Login - 
func Login(region string) {
  output, err := commands.Run(str.Concat("aws ecr get-login --no-include-email --region " + region), "")
  errors.LogIfError(err)
  commands.Run(output, "")
}
