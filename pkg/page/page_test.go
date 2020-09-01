package page

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetByRoute(t *testing.T) {
	files.Root = "../files/"
	GetByRoute("test", "test_page", false)
}

func TestGet(t *testing.T) {
	files.Root = "../files/"
	Get("test", 1, false)
}

func TestSaveData(t *testing.T) {
	files.Root = "../files/"
	pages = []Info{Info{Title: "test"}}
	SaveData("test")
	data, err := ioutil.ReadFile(files.Root + "test" + pagesFile)
	err = json.Unmarshal(data, &pages)
	if err != nil {
		t.Error("Failed to decode json data")
	}
	title := pages[0].Title
	if title != "test" {
		t.Errorf("Got %s want test", title)
	}
}

func TestLoadData(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	title := pages[0].Title
	if title != "test" {
		t.Errorf("Got %s want test after first running TestSaveData", title)
	}
}

func TestSaveContent(t *testing.T) {
	files.Root = "../files/"
	SaveContent("test", 9999, "Test")
	fn := files.Root + "test/pages/9999.html"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("Failed to save file")
	}
	content := string(data)
	if content != "Test" {
		t.Errorf("Got %s want Test", content)
	}
	os.Remove(fn)
}

func TestGeneratePages(t *testing.T) {
	files.Root = "../files/"
	// Set the root empty page
	pages = []Info{Info{}}
	id := 0
	for a := 0; a < 4; a++ {
		id++
		aid := id
		addPage(fmt.Sprintf("%d", a), aid, 0, "home")
		for b := 0; b < 4; b++ {
			id++
			bid := id
			addPage(fmt.Sprintf("%d-%d", a, b), bid, aid, "page")
			for c := 0; c < 4; c++ {
				id++
				cid := id
				addPage(fmt.Sprintf("%d-%d-%d", a, b, c), cid, bid, "post")
			}
		}
	}
	SaveData("test")
}

func addPage(route string, id int, parent int, template string) {
	title := "Page " + route
	if route == "0" {
		route = ""
	}
	pages = append(pages, Info{ID: id, Parent: parent, Title: title, Content: title, Template: template, Route: "/" + route})
	SaveContent("test", id, title)
}
