package util

import (
  "strings"
)

func IsArn(arn string) bool {
  return strings.HasPrefix(arn, "arn:")
}
