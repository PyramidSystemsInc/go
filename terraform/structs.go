package terraform

// Provider - Values used to create a Terraform AWS Provider
type Provider struct {
  ProjectName     string
  Region          string
  AWSVersion      string
  TemplateVersion string
}
