package main

import (
	"fmt"
	"gocms/scripts/util"
	"io/ioutil"
	"path/filepath"
)

func main() {
	path := util.CheckArgs()

	// Create folders
	util.CreateFolder(path, "data")
	util.CreateFolder(path, "templates")
	util.CreateFolder(path, "pages")

	// Copy files
	files, err := ioutil.ReadDir("build")
	util.Check(err)
	for _, f := range files {
		fn := f.Name()
		if filepath.Ext(fn) != ".md" {
			fmt.Printf("Copied: %s\n", fn)
			util.CopyFile(filepath.Join("build", fn), filepath.Join(path, fn))
		}
	}
}
