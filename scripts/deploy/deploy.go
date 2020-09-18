// The first command argument should be the path to the domain folder and it's parent directory will be where the App is installed
package main

import (
	"fmt"
	"gocms/scripts/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	path := util.CheckArgs()
	appPath := filepath.Dir(path)
	_, err := os.Stat(appPath)
	if err != nil {
		util.Check(fmt.Errorf("The given App Path %s does not exist", appPath))
	}

	// Create folders
	util.CreateFolder(path, "")
	util.CreateFolder(path, "data")
	util.CreateFolder(path, "templates")
	util.CreateFolder(path, "pages")
	util.CreateFolder(path, "static")

	// Copy files
	files, err := ioutil.ReadDir("build")
	util.Check(err)
	for _, f := range files {
		fn := f.Name()
		if filepath.Ext(fn) != ".md" {
			fmt.Printf("Copied: %s\n", fn)
			destFile := filepath.Join(appPath, fn)
			util.CopyFile(filepath.Join("build", fn), destFile)
			os.Chmod(destFile, 0755)
		}
	}
}
