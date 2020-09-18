// Copy test files to website folders
package main

import (
	"fmt"
	"gocms/scripts/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

var websitePath string

func main() {
	path := util.CheckArgs()

	util.CreateFolder(path, "js")
	util.CreateFolder(path, "css")

	// Remove existing page files
	destPath := filepath.Join(path, "pages")
	removeFiles(destPath)

	// Copy page content files
	srcPath := "./pkg/files/test/pages"
	count := copyFiles(srcPath, destPath, ".html")
	fmt.Printf("Copied %d page content files.\n", count)

	// Remove existing template files
	destPath = filepath.Join(path, "templates")
	removeFiles(destPath)

	// Copy template files
	srcPath = "./pkg/files/test/templates"
	count = copyFiles(srcPath, destPath, ".html")
	fmt.Printf("Copied %d template files.\n", count)

	// Remove existing css files
	destPath = filepath.Join(path, "static", "css")
	removeFiles(destPath)

	// Copy css files
	srcPath = "./pkg/files/test/styles"
	count = copyFiles(srcPath, destPath, ".css")
	fmt.Printf("Copied %d css files.\n", count)

	// Remove existing js files
	destPath = filepath.Join(path, "static", "js")
	removeFiles(destPath)

	// Copy js files
	srcPath = "./pkg/files/test/scripts"
	count = copyFiles(srcPath, destPath, ".js")
	fmt.Printf("Copied %d script files.\n", count)

	// Remove existing data files
	destPath = filepath.Join(path, "data")
	removeFiles(destPath)

	// Copy data files
	srcPath = "./pkg/files/test/data"
	count = copyFiles(srcPath, destPath, ".json")
	fmt.Printf("Copied %d data files.\n", count)

}

func removeFiles(path string) {
	os.RemoveAll(path)
	err := os.Mkdir(path, 0755)
	util.Check(err)
}

func copyFiles(srcPath string, destPath string, ext string) int {
	files, err := ioutil.ReadDir(srcPath)
	util.Check(err)
	count := 0
	for _, f := range files {
		fn := f.Name()
		if filepath.Ext(fn) == ext {
			util.CopyFile(filepath.Join(srcPath, fn), filepath.Join(destPath, fn))
			count++
		}
	}
	return count
}
