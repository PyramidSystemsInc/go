package files

import (
  "io/ioutil"
  "net/http"
  "os"
  "strings"
  "text/template"
  "github.com/PyramidSystemsInc/go/directories"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
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

func Read(filePath string) []byte {
  data, err := ioutil.ReadFile(filePath)
  errors.LogIfError(err)
  return data
}

// TODO: Do some regex checking on valid values of fullPath
func Download(url string, fullPath string) {
  resp, err := http.Get(url)
  errors.LogIfError(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  errors.LogIfError(err)
  Write(fullPath, body)
}

func Write(fullPath string, data []byte) {
  ioutil.WriteFile(fullPath, data, 0644)
}

func ChangePermissions(fullPath string, permissions int) {
  if strings.Index(fullPath, ".") == 0 {
    fullPath = str.Concat(directories.GetWorking(), fullPath[1:len(fullPath)])
  }
  err := os.Chmod(fullPath, os.FileMode(permissions))
  errors.LogIfError(err)
}
