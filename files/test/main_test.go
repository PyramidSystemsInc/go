package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/PyramidSystemsInc/go/files"
)

func TestFileStringReplace(t *testing.T) {
	file := "./testFile.txt"

	//create file with contents of "Hello world!"
	err := ioutil.WriteFile(file, []byte("Hello world!"), 0644)
	if err != nil {
		t.Error()
	}

	defer os.Remove(file)

	//replace "world" with "Mom"
	files.FileStringReplace("./testFile.txt", "world", "Mom")

	//check that file contents is "Hello Mom!"
	f, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error()
	}

	if string(f) != "Hello Mom!" {
		t.Error()
	}

}

func TestDirectoryFileStringReplace(t *testing.T) {
	file1 := "./testFile_1.txt"
	file2 := "./testFile_2.txt"

	//create file with contents of "Hello world!"
	err := ioutil.WriteFile(file1, []byte("Hello world!"), 0644)
	if err != nil {
		t.Error()
	}

	defer os.Remove(file1)

	err = ioutil.WriteFile(file2, []byte("Hello world!"), 0644)
	if err != nil {
		t.Error()
	}

	defer os.Remove(file2)

	files.DirectoryFileStringReplace(".", "testFile", "world", "Mom")

	result1, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Error()
	}

	if string(result1) != "Hello Mom!" {
		t.Error()
	}

	result2, err := ioutil.ReadFile(file2)
	if err != nil {
		t.Error()
	}

	if string(result2) != "Hello Mom!" {
		t.Error()
	}
}
