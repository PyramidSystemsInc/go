package files

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/PyramidSystemsInc/go/directories"
	"github.com/PyramidSystemsInc/go/errors"
	"github.com/PyramidSystemsInc/go/str"
)

func Delete(filePath string) error {
	err := os.Remove(filePath)
	return err
}

// Simply creates a file given a full path (including file name and extension)
func CreateBlank(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	errors.LogIfError(err)
	return file
}

func FindUpTree(fileName string) string {
	workingDirectory := directories.GetWorking()
	directory := workingDirectory
	for {
		filePath := path.Join(directory, fileName)
		if Exists(filePath) {
			return directory
		}
		if directory == "/" {
			break
		}
		directory = path.Join(directory, "..")
	}
	return ""
}

func Exists(filePath string) bool {
	filePath = path.Clean(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Ensure the full path exists
func EnsurePath(filePath string) {
	os.MkdirAll(filePath, 0755)
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
func Download(url string, fullPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	Write(fullPath, body)
	return err
}

func Write(fullPath string, data []byte) {
	ioutil.WriteFile(fullPath, data, 0644)
}

func Prepend(filePath string, data []byte) {
	content := Read(filePath)
	newContent := str.Concat(string(data), string(content))
	Write(filePath, []byte(newContent))
}

func Append(filePath string, data []byte) {
	content := Read(filePath)
	newContent := str.Concat(string(content), string(data))
	Write(filePath, []byte(newContent))
}

func AppendBelow(filePath string, markerLine string, data string) {
	var newFile string
	file, err := os.Open(filePath)
	if err != nil {
		errors.LogAndQuit(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newFile += scanner.Text() + "\n"
		if scanner.Text() == markerLine {
			newFile += data + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		errors.LogAndQuit(err.Error())
	}
	Write(filePath, []byte(newFile))
}

func ChangePermissions(fullPath string, permissions int) {
	if strings.Index(fullPath, ".") == 0 {
		fullPath = str.Concat(directories.GetWorking(), fullPath[1:len(fullPath)])
	}
	err := os.Chmod(fullPath, os.FileMode(permissions))
	errors.LogIfError(err)
}
