package util

import (
  "strings"
)

func IsArn(possibleArn string) bool {
  arnPrefix := "arn:aws:"
  return strings.HasPrefix(possibleArn, arnPrefix)
}
