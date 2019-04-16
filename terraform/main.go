package terraform

import (
  "github.com/PyramidSystemsInc/go/commands"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

// Attempts to get the Terraform version to demonstrate Terraform is installed and accessible
// If Terraform is not installed or accessible the execution of the program is stopped
func VerifyInstallation() {
  _, err := commands.Run("terraform version", "")
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Checking the Terraform version failed with the following error: ", err.Error()))
  }
}


// Initializes the terraform directory checks for *.tf files and processes them
func Initialize(directoryToRunFrom string) string {
  output, err := commands.Run("terraform init -input=false", directoryToRunFrom)
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Initializing Terraform failed with the following error: ", err.Error()))
  }
  return output
}

