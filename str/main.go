package str

import (
  "strings"
)

// Concat - Concatinates an arbitrary number of strings
func Concat(firstString string, moreStrings... string) string {
  var completeString string
  completeString += trimCarriageReturnSuffix(firstString)
  for _, thisString := range moreStrings {
    completeString += trimCarriageReturnSuffix(thisString)
  }
  return completeString
}

// IsAllLowercaseCharacters - Tests whether or not a string contains only lowercase, alphabetic characters
func IsAllLowercaseCharacters(myString string) bool {
  return strings.IndexFunc(myString, isNotLowerCaseCharacter) == -1
}

// RunForEachLine - Runs a provided function (with a single string parameter input and no output) for each line of a multiline string
func RunForEachLine(function func(string), multilineString string) {
  newLine := "\n"
  for _, line := range strings.Split(strings.TrimSuffix(multilineString, newLine), newLine) {
    function(line)
  }
}

func isNotLowerCaseCharacter(character rune) bool {
  return character < 'a' || character > 'z'
}

func trimCarriageReturnSuffix(myString string) string {
  return strings.TrimSuffix(strings.TrimSuffix(myString, "\n"), "\r")
}
