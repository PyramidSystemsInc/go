package str

import (
  "strings"
)

func Concat(firstString string, moreStrings... string) string {
  var completeString string
  completeString += trimCarriageReturnSuffix(firstString)
  for _, thisString := range moreStrings {
    completeString += trimCarriageReturnSuffix(thisString)
  }
  return completeString
}

func trimCarriageReturnSuffix(myString string) string {
  return strings.TrimSuffix(strings.TrimSuffix(myString, "\n"), "\r")
}

func IsAllLowercaseCharacters(myString string) bool {
  return strings.IndexFunc(myString, isNotLowerCaseCharacter) == -1
}

func isNotLowerCaseCharacter(character rune) bool {
  return character < 'a' || character > 'z'
}
