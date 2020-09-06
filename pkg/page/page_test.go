package page

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGetByRoute(t *testing.T) {
	files.Root = "../files/"
	GetByRoute("test", "test_page", false)
}

func TestGetByID(t *testing.T) {
	files.Root = "../files/"
	p := GetByID("test", 2, false)
	title := p.Title
	want := "Page 0-0"
	if title != want {
		t.Errorf("Want %s got %s", want, title)
	}
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
	pubDate := time.Now()
	// Set the root empty page
	pages = []Info{Info{Parent: -1}}
	id := 0
	for a := 0; a < 4; a++ {
		id++
		aid := id
		addPage("Page ", fmt.Sprintf("%d", a), aid, 0, "home", 0, pubDate)
		for b := 0; b < 4; b++ {
			id++
			bid := id
			addPage("Page ", fmt.Sprintf("%d-%d", a, b), bid, aid, "page", 0, pubDate)
			for c := 0; c < 4; c++ {
				id++
				cid := id
				addPage("Page ", fmt.Sprintf("%d-%d-%d", a, b, c), cid, bid, "past", 1, pubDate)
			}
		}
	}
	SaveData("test")
}

func TestGeneratePosts(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	id := 300
	for a := 0; a < 4; a++ {
		id++
		aid := id
		addPage("Post ", fmt.Sprintf("blog/%d", a), aid, 0, "post", 0, getTime(2000+a*5, time.January, a*2))
		for b := 0; b < 4; b++ {
			id++
			bid := id
			addPage("Post ", fmt.Sprintf("blog/%d-%d", a, b), bid, aid, "post", 0, getTime(2000+a*5, time.Month(2+b), 1))
			for c := 0; c < 4; c++ {
				id++
				cid := id
				addPage("Post ", fmt.Sprintf("blog/%d-%d-%d", a, b, c), cid, bid, "post", 0, getTime(2000+a*5, time.Month(2+b), 1+c*2))
			}
		}
	}
	SaveData("test")
}

func getTime(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 23, 0, 0, 0, time.UTC)
}

func TestGenerateCategoryPages(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	pubDate := time.Now()
	// Set the root category page
	id := 600
	addPage("Categories ", "cats", id, 0, "category", 0, pubDate)
	for a := 0; a < 4; a++ {
		id++
		aid := id
		addPage("Category ", fmt.Sprintf("cat-%d", aid), aid, 999, "category", 999, pubDate)
		for b := 0; b < 4; b++ {
			id++
			bid := id
			addPage("Category ", fmt.Sprintf("cat-%d-%d", aid, bid), bid, aid, "category", aid, pubDate)
			for c := 0; c < 4; c++ {
				id++
				cid := id
				addPage("Category ", fmt.Sprintf("cat-%d-%d-%d", aid, bid, cid), cid, bid, "category", bid, pubDate)
			}
		}
	}
	SaveData("test")
}

func TestGetPages(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	list := GetPages(0, Published)
	if len(list) != 4 {
		t.Errorf("Got length of %v want 4", len(list))
	}
}

func addPage(title string, route string, id int, parent int, template string, category int, pubDate time.Time) {
	title += route
	if route == "0" {
		route = ""
	}
	menu := ""
	if id < 5 {
		menu = "side"
	}
	Add(Info{
		ID:       id,
		Parent:   parent,
		Title:    title,
		Content:  title,
		Template: template,
		Route:    "/" + route,
		Status:   Published,
		Category: category,
		Menu:     menu,
		PubDate:  pubDate,
	})
	SaveContent("test", id, title)
}

func TestGetRecentPosts(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	want := 4
	posts := GetRecentPosts(want)
	got := len(posts)
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}
