package page

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Status type
type Status int

// Type type
type Type int

// Draft or Published status constants
const (
	Draft Status = iota
	Published
)

// Post or Page type constants
const (
	Page Type = iota
	Post
)

// Info type
type Info struct {
	ID          int
	Parent      int
	Depth       int
	Title       string
	Description string
	Content     string
	Route       string
	Author      string
	Status      Status
	Category    int
	Type        Type
	Template    string
	Timestamp   int
}

var pages []Info

const pagesFile = "/data/pages.json"

// Get by ID a page or post
func Get(domain string, id int, getContent bool) Info {
	return find(domain, id, "", getContent)
}

// GetByRoute a page or post
func GetByRoute(domain string, route string, getContent bool) Info {
	return find(domain, -1, route, getContent)
}

func find(domain string, id int, route string, getContent bool) Info {
	var page Info
	if pages == nil {
		LoadData(domain)
	}
	for _, p := range pages {
		if p.ID == id || p.Route == route {
			if getContent {
				p.Content = LoadContent(domain, p.ID)
			}
			page = p
			break
		}
	}
	return page //Info{Title: "page title", Route: "/page-route", Timestamp: 3456}
}

// LoadContent return the contents of the page content file
func LoadContent(domain string, id int) string {
	content, _ := ioutil.ReadFile(fmt.Sprintf("%s/pages/%d.html", domain, id))
	return string(content)
}

// LoadData loads the data from the pages data file
func LoadData(domain string) {
	data, err := ioutil.ReadFile(domain + pagesFile)
	if err != nil {
		pages = []Info{Info{}}
		b, _ := json.Marshal(pages)
		data = b
		err := ioutil.WriteFile(domain+pagesFile, b, 0660)
		if err != nil {
			panic(err)
		}
	}
	err = json.Unmarshal(data, &pages)
	if err != nil {
		panic(err)
	}
}

// Save a page or post
func Save(info Info, saveContent bool) int {
	return 1
}

// List - Get a list of all pages
func List() []Info {
	return []Info{Info{}, Info{}}
}
