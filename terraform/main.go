package terraform

import (
  "fmt"
  "time"

  "github.com/PyramidSystemsInc/go/commands"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/files"
  "github.com/PyramidSystemsInc/go/str"
)

// Apply - Creates resources detailed in the tfplan file (created using the `terraform plan` command
func Apply(directoryToRunFrom string) string {
  defer timeTrack(time.Now(), "Terraform apply")
  output, err := commands.Run("terraform apply -input=false tfplan", directoryToRunFrom)
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Applying the Terraform plan failed with the following error: ", err.Error()))
  }
  return output
}

// Destroy - Destroys all resources managed by Terraform
func Destroy(directoryToRunFrom string) string {
  defer timeTrack(time.Now(), "Terraform destroy")
  if files.Exists(str.Concat(directoryToRunFrom, "/.terraform")) {
    output, err := commands.Run("terraform destroy -auto-approve", directoryToRunFrom)
    if err != nil {
      errors.LogAndQuit(str.Concat("ERROR: Terraform destroy failed with the following error: ", err.Error()))
    }
    return output
  } else {
    return str.Concat("No Terraform resources to destroy in ", directoryToRunFrom)
  }
}

// Initialize - Initializes the terraform directory, checks for *.tf files, and processes them
func Initialize(directoryToRunFrom string) string {
  output, err := commands.Run("terraform init -input=false", directoryToRunFrom)
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Initializing Terraform failed with the following error: ", err.Error()))
  }
  return output
}

// Plan - Creates a tfplan file with a detailed specification of what Terraform would create given the set of *.tf files
func Plan(directoryToRunFrom string, cfg map[string]string) string {
  var variables string
  for key, value := range cfg {
    variables = str.Concat(variables, "-var ", key, "=", value, " ")
  }
  planCommand := str.Concat("terraform plan ", variables, "-out tfplan")
  output, err := commands.Run(planCommand, directoryToRunFrom)
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Initializing Terraform failed with the following error: ", err.Error()))
  }
  return output
}

// VerifyInstallation - Attempts to get the Terraform version to demonstrate Terraform is installed and accessible. If Terraform is not installed or accessible the execution of the program is stopped
func VerifyInstallation() {
  _, err := commands.Run("terraform version", "")
  if err != nil {
    errors.LogAndQuit(str.Concat("ERROR: Checking the Terraform version failed with the following error: ", err.Error()))
  }
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  fmt.Sprintf("%s took %s", name, elapsed)
}
