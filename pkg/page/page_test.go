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

func TestCountPagesInCategory(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	got := CountPagesInCategory(601)
	want := 3
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGeneratePages(t *testing.T) {
	files.Root = "../files/"
	pubDate := time.Now()
	// Set the root empty page
	pages = []Info{Info{Parent: -1}}
	id := 0
	title := "Home Page"
	template := "home"
	cat := 0
	for a := 0; a < 4; a++ {
		id++
		aid := id
		addPage(title, fmt.Sprintf("%d", a), aid, 0, template, cat, pubDate)
		title = "Page "
		template = "page"
		cat = 601
		for b := 0; b < 4; b++ {
			id++
			bid := id
			addPage(title, fmt.Sprintf("%d-%d", a, b), bid, aid, "page", 602, pubDate)
			for c := 0; c < 4; c++ {
				id++
				cid := id
				addPage(title, fmt.Sprintf("%d-%d-%d", a, b, c), cid, bid, "page", 603, pubDate)
			}
		}
	}
	id++
	addPage("Archives", "archive", id, 0, "archive", 0, time.Now())
	id++
	addPage("Blog", "blog", id, 0, "blog", 0, time.Now())
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
		addPage("Category ", fmt.Sprintf("cat-%d", aid), aid, 600, "category", 600, pubDate)
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
	list := GetPages(1, Published, "page")
	got := len(list)
	want := 4
	if got != want {
		t.Errorf("Got length of %d want %d", got, want)
	}
}

func addPage(title string, route string, id int, parent int, template string, category int, pubDate time.Time) {
	if template == "post" || template == "page" || id > 600 && template == "category" {
		title += route
	}
	if route == "0" {
		route = ""
	}
	menu := "-"
	switch id {
	case 1, 85, 86, 600:
		menu = "top"
	case 2, 3, 4, 5:
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
		Menus:    []string{menu},
		PubDate:  pubDate,
	})

	content := "Content for " + title
	if id == 384 {
		content = "before#MORE#after"
	}

	SaveContent("test", id, content)
}

func TestGetRecentPosts(t *testing.T) {
	files.Root = "../files/"
	LoadData("test")
	want := 10
	posts := GetRecentPosts(want, false, "")
	got := len(posts)
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}
