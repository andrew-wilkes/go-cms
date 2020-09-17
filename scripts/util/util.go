package util

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// CheckArgs func
func CheckArgs() string {
	if len(os.Args) < 2 {
		Check(errors.New("Missing installation path"))
	}
	path := os.Args[1]
	return path
}

// CreateFolder if it doesn't already exist
func CreateFolder(path string, name string) {
	if name != "" {
		path = filepath.Join(path, name)
	}
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0755)
		Check(err)
	}
}

// CopyFile func
func CopyFile(src string, dest string) {
	original, err := os.Open(src)
	Check(err)
	defer original.Close()

	new, err := os.Create(dest)
	Check(err)
	defer new.Close()

	_, err = io.Copy(new, original)
	Check(err)
}

// Check for error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
