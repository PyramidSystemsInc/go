package directories

import (
	"os"
	"os/user"
	"path"

	"github.com/PyramidSystemsInc/go/errors"
)

// Change - Changes the working directory to the directory specified (like `cd`)
func Change(directory string) error {
	return os.Chdir(directory)
}

// Create - Creates a directory (if it does not already exist)
func Create(directory string) {
	directory = path.Clean(directory)
	err := os.MkdirAll(directory, os.ModePerm)
	errors.QuitIfError(err)
}

// GetHome - Returns the home directory of the user running the program
func GetHome() string {
	user, err := user.Current()
	errors.QuitIfError(err)
	return user.HomeDir
}

// GetWorking - Returns the working directory
func GetWorking() string {
	workingDirectory, err := os.Getwd()
	errors.LogIfError(err)
	return path.Clean(workingDirectory)
}

// Exists - returns whether the given file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	errors.LogIfError(err)
	return true
}
