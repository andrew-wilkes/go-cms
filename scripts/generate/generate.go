package main

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"time"
)

var domain string

func main() {
	Run()
}

// Run the page data generator functions
func Run() {
	files.Root = "pkg/files/"
	domain = "test"
	generatePages()
	generatePosts()
	generateCategoryPages()
}

func generatePages() {
	pubDate := time.Now()
	// Set the root empty page
	b, _ := json.Marshal([]page.Info{page.Info{Parent: -1}})
	page.SaveRawData(domain, b)
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
	page.SaveData(domain)
}

func generatePosts() {
	page.LoadData(domain)
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
	page.SaveData(domain)
}

func generateCategoryPages() {
	page.LoadData(domain)
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
	page.SaveData(domain)
}

func getTime(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 23, 0, 0, 0, time.UTC)
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

	page.Add(page.Info{
		ID:       id,
		Parent:   parent,
		Title:    title,
		Content:  title,
		Template: template,
		Route:    "/" + route,
		Status:   page.Published,
		Category: category,
		Menus:    []string{menu},
		PubDate:  pubDate,
	})

	content := "Content for " + title
	if id == 384 {
		content = "before#MORE#after"
	}

	page.SaveContent("test", id, content)
}
