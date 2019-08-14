package files

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/PyramidSystemsInc/go/directories"
	"github.com/PyramidSystemsInc/go/errors"
	"github.com/PyramidSystemsInc/go/logger"
	"github.com/PyramidSystemsInc/go/str"
	"github.com/gobuffalo/packr"
)

// Delete - Deletes a single file
func Delete(filePath string) error {
	err := os.Remove(filePath)
	return err
}

// CreateBlank - Creates a file given a full path (including file name and extension) that has no contents
func CreateBlank(filePath string) *os.File {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	errors.LogIfError(err)
	return file
}

// CreateFromBinary creates a file populates it with the contents provided
func CreateFromBinary(filePath string, fileContents []byte) {
	file, err := os.Create(filePath)
	errors.QuitIfError(err)

	_, err = file.Write(fileContents)
	errors.QuitIfError(err)

	file.Close()
}

// FindUpTree - Looks recursively in parent directories until a certain file name is found and returns the path to that file
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

// Exists - Checks if a file exists
func Exists(filePath string) bool {
	filePath = path.Clean(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// EnsurePath - Creates a path if it does not already exist
func EnsurePath(filePath string) {
	os.MkdirAll(filePath, 0755)
}

// Read - Returns the contents of a file
func Read(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	errors.LogIfError(err)
	return data
}

// TODO: Do some regex checking on valid values of fullPath
// Download - Downloads a file from a URL to a given path on the local filesystem
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

// Write - Writes (or overwrites) a file given a full path and contents
func Write(fullPath string, data []byte) {
	ioutil.WriteFile(fullPath, data, 0644)
}

// Prepend - Adds content to the top of a file
func Prepend(filePath string, data []byte) {
	content := Read(filePath)
	newContent := str.Concat(string(data), string(content))
	Write(filePath, []byte(newContent))
}

// Append - Adds content to the bottom of a file
func Append(filePath string, data []byte) {
	content := Read(filePath)
	newContent := str.Concat(string(content), string(data))
	Write(filePath, []byte(newContent))
}

// AppendBelow - Adds content to the middle of a file directly below a line matching `markerLine`
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

// ChangePermissions - Changes the permissions of a file
func ChangePermissions(fullPath string, permissions int) {
	if strings.Index(fullPath, ".") == 0 {
		fullPath = str.Concat(directories.GetWorking(), fullPath[1:len(fullPath)])
	}
	err := os.Chmod(fullPath, os.FileMode(permissions))
	errors.LogIfError(err)
}

// TemplateOptions holds the various pieces that might be needed to output a single file
// Not all options are required
type TemplateOptions struct {
	TargetDirectory string
	Box             packr.Box
	TemplatePath    string
	Config          map[string]string
	FileRenames     map[string]string
}

// CreateTemplatedFiles creates all files specified by options passed in
func CreateTemplatedFiles(options TemplateOptions) error {
	// sanity check inputs
	box := options.Box
	if options.Box.Path == "" {
		return errors.New("No box passed in")
	}

	for _, templatePath := range box.List() {
		options.TemplatePath = templatePath
		CreateTemplatedFile(options)
	}
	return nil
}

// CreateTemplatedFile creates a single file, specified by the passed-in TemplateOptions
func CreateTemplatedFile(options TemplateOptions) error {
	// constants
	// note that template.ts is not actually binary, just needs to be handled as pass-through
	binaries := []string{".zip", ".png", ".jpg", ".ico", ".jar", ".ear", ".war", ".template.ts"}

	// sanity check inputs
	box := options.Box
	if options.Box.Path == "" {
		return errors.New("No box passed in")
	}
	if options.TemplatePath == "" {
		return errors.New("No templatePath passed in")
	}

	// target is the templatePath
	fullPath := options.TemplatePath
	// unless we have a rename for it
	newName := options.FileRenames[options.TemplatePath]
	if newName != "" {
		fullPath = newName
	}
	// see if the filePath needs template eval
	if strings.Contains(fullPath, "{{") {
		t := template.Must(template.New("t").Parse(fullPath))
		buf := new(strings.Builder)
		t.Execute(buf, options.Config)
		fullPath = buf.String()
	}
	// maybe put it in a subdirectory
	if options.TargetDirectory != "" {
		fullPath = filepath.Join(options.TargetDirectory, fullPath)
	}
	EnsurePath(filepath.Dir(fullPath))

	// check if we are dealing with a binary file or one whose contents need template replacing
	isBinary := false
	for _, each := range binaries {
		isBinary = strings.HasSuffix(options.TemplatePath, each)
		if isBinary {
			break
		}
	}
	if isBinary {
		fileContent, err := box.Find(options.TemplatePath)
		errors.QuitIfError(err)
		CreateFileWithContent(fullPath, fileContent)
	} else {
		template, err := box.FindString(options.TemplatePath)
		errors.QuitIfError(err)
		CreateFromTemplate(fullPath, template, options.Config)
	}
	logger.Info("Created " + fullPath)
	return nil
}

// CreateFromTemplate - Creates a file and populates it with a given template
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

// CreateFileWithContent creates a file and puts the data provided into it
func CreateFileWithContent(filePath string, content []byte) {
	file, err := os.Create(filePath)
	errors.QuitIfError(err)
	_, err = file.Write(content)
	errors.QuitIfError(err)
	file.Close()
}

// FileStringReplace opens a given file and peforms a search and replace of its contents and saves the file.
func FileStringReplace(file, find, replace string) {
	//copy file contents into variable
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	//check the regex is valid
	re := regexp.MustCompile(find)

	//grep variable for string
	//replace all occurences of string in variable
	s := re.ReplaceAllString(string(data), replace)

	//overwrite file with new contents
	w := ioutil.WriteFile(file, []byte(s), 0666)
	if w != nil {
		log.Fatalf("unable to write to file: %s\n", w)
	}
}

// DirectoryFileStringReplace reads the target directory for files matching the filex regex pattern.
// It then performs a search and replace in the matching files.
func DirectoryFileStringReplace(dir, filex, search, replace string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("unable to read file")
		log.Fatal(err)
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		//check the regex is valid
		re := regexp.MustCompile(filex)

		if re.MatchString(file.Name()) {
			FileStringReplace(path, search, replace)
		}
	}
}
