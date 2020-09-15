package page

import (
	"encoding/json"
	"gocms/pkg/files"
	"io/ioutil"
	"os"
	"testing"
)

var domain string

func TestSaveData(t *testing.T) {
	setup(t)
	backup, err := ioutil.ReadFile(files.Root + domain + pagesFile)
	pages = []Info{Info{Title: domain}}
	SaveData(domain)
	data, err := ioutil.ReadFile(files.Root + domain + pagesFile)
	err = json.Unmarshal(data, &pages)
	if err != nil {
		t.Error("Failed to decode json data")
	}
	title := pages[0].Title
	if title != domain {
		t.Errorf("Got %s want test", title)
	}
	err = ioutil.WriteFile(files.Root+domain+pagesFile, backup, 0660)
	if err != nil {
		panic(err)
	}
}

func TestLoadData(t *testing.T) {
	setup(t)
	title := pages[0].Title
	want := "Home Page"
	if title != want {
		t.Errorf("Got %s want %s", title, want)
	}
}

func TestCountPagesInCategory(t *testing.T) {
	setup(t)
	got := CountPagesInCategory(602)
	want := 16
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGetPages(t *testing.T) {
	setup(t)
	list := GetPages(1, Published, "page")
	got := len(list)
	want := 4
	if got != want {
		t.Errorf("Got length of %d want %d", got, want)
	}
}

func TestGetRecentPosts(t *testing.T) {
	setup(t)
	want := 10
	posts := GetRecentPosts(want, false, "")
	got := len(posts)
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGetByRoute(t *testing.T) {
	setup(t)
	GetByRoute(domain, "test_page", false)
}

func TestGetByID(t *testing.T) {
	setup(t)
	p := GetByID(domain, 2, false)
	title := p.Title
	want := "Page 0-0"
	if title != want {
		t.Errorf("Want %s got %s", want, title)
	}
}

func TestSaveContent(t *testing.T) {
	setup(t)
	SaveContent(domain, 9999, domain)
	fn := files.Root + "test/pages/9999.html"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("Failed to save file")
	}
	content := string(data)
	if content != domain {
		t.Errorf("Got %s want Test", content)
	}
	os.Remove(fn)
}

func setup(t *testing.T) {
	files.Root = "../files/"
	domain = "test"
	LoadData(domain)
	if len(pages) < 100 {
		t.Errorf(`Need to generate pages data. Run: "Make g" in terminal`)
	}
}
