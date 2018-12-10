package files

import (
  "os"
  "text/template"
  "github.com/PyramidSystemsInc/go/errors"
)

// Simply creates a file given a full path (including file name and extension)
func CreateBlank(filePath string) *os.File {
  file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
  errors.LogIfError(err)
  return file
}

// Creates a file and populates it with a given template
// If a template features the following syntax: {{.mapKey}}, the value of
//   'mapKey' in the config variable will be inserted
func CreateFromTemplate(filePath string, pattern string, config map[string]string) {
  t := template.Must(template.New("t").Parse(pattern))
  file, err := os.Create(filePath)
  errors.QuitIfError(err)
  err = t.Execute(file, config)
  errors.QuitIfError(err)
  file.Close()
}
