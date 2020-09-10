package page

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"io/ioutil"
	"time"
)

// Status type
type Status int

// Draft or Published status constants
const (
	Draft Status = iota
	Published
)

// Info type
type Info struct {
	ID          int       // This is a reference number for the page data
	Parent      int       // The ID of this page's parent in the tree
	Depth       int       // The indentation level for this page in the Dashboard tree
	Title       string    // As used in the HTML title tag and link text
	Description string    // As used in the HTML description tag and link title text (optional)
	Content     string    // Used to store the HTML content for the page when it needs to be passed around in the Dashboard
	Route       string    // The virtual URI route to match with the webpage being requested
	Author      string    // Maybe useful for displaying the author of the content
	Status      Status    // A flag for the published or draft status of the page
	Category    int       // Each page may be associated with one category (do we need tags as well?)
	Menus       []string  // The lower-case names of menus that this page is associated with such as: side, top, footer
	Template    string    // Base name of the template file to use such as "page" for the page.html template
	PubDate     time.Time // Time when the page was first saved
	UpdateDate  time.Time // Time of the latest re-save
}

// GetByID by ID a page or post
func GetByID(domain string, id int, getContent bool) Info {
	return find(domain, id, "-", getContent)
}

// GetByRoute a page or post
func GetByRoute(domain string, route string, getContent bool) Info {
	return find(domain, -1, route, getContent)
}

var pages []Info

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
	return page
}

// LoadContent return the contents of the page content file
func LoadContent(domain string, id int) string {
	content, _ := ioutil.ReadFile(fmt.Sprintf("%s%s/pages/%d.html", files.Root, domain, id))
	return string(content)
}

// SaveContent saves page content to a file
func SaveContent(domain string, id int, content string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s%s/pages/%d.html", files.Root, domain, id), []byte(content), 0660)
	if err != nil {
		panic(err)
	}
}

const pagesFile = "/data/pages.json"

// LoadData loads the data from the pages data file
func LoadData(domain string) {
	data, err := ioutil.ReadFile(files.Root + domain + pagesFile)
	if err != nil {
		pages = []Info{Info{}}
		SaveData(domain)
	} else {
		err = json.Unmarshal(data, &pages)
		if err != nil {
			panic(err)
		}
	}
}

// SaveData saves the pages data to a file
func SaveData(domain string) {
	b, _ := json.Marshal(pages)
	err := ioutil.WriteFile(files.Root+domain+pagesFile, b, 0660)
	if err != nil {
		panic(err)
	}
}

// Save a page or post
func Save(domain string, info Info, saveContent bool) int {
	if saveContent {
		SaveContent(domain, info.ID, info.Content)
	}
	return 1
}

// GetPages returns a slice of pages data
func GetPages(parent int, status Status, template string) []Info {
	list := []Info{}
	for _, p := range pages {
		if p.Parent == parent && p.Status == status && p.Template == template {
			list = append(list, p)
		}
	}
	return list
}

// GetAllPages return all the page data
func GetAllPages() []Info {
	return pages
}

// GetCategoryPages returns a slice of category pages in a category
func GetCategoryPages(id int, status Status) []Info {
	list := []Info{}
	for _, p := range pages {
		if p.Category == id && p.Status == status && p.Template == "category" {
			list = append(list, p)
		}
	}
	return list
}

// GetPagesInCategory returns a slice of pages and posts in a category
func GetPagesInCategory(id int) []Info {
	list := []Info{}
	for _, p := range pages {
		if p.Category == id && p.Status == Published && p.Template != "category" {
			list = append(list, p)
		}
	}
	return list
}

// CountPagesInCategory returns the count of pages in a category
func CountPagesInCategory(id int) int {
	n := 0
	for _, p := range pages {
		if p.Category == id && p.Status == Published && p.Template != "category" {
			n++
		}
	}
	return n
}

// GetPagesInMenu returns a slice of pages and posts in a menu
func GetPagesInMenu(menu string) []Info {
	list := []Info{}
	for _, p := range pages {
		if p.Status == Published {
			for _, m := range p.Menus {
				if menu == m {
					list = append(list, p)
					break
				}
			}

		}
	}
	return list
}

// Add new page data
func Add(newPage Info) {
	pages = append(pages, newPage)
}

// GetRecentPosts returns a slice of the posts with highest index value
func GetRecentPosts(n int, getContent bool, domain string) []Info {
	var posts = make([]Info, n)
	count := 0
	for i := len(pages) - 1; i >= 0; i-- {
		p := pages[i]
		if p.Status == Published && p.Template == "post" {
			posts[count] = p
			if getContent {
				p.Content = LoadContent(domain, p.ID)
			}
			count++
			if count == n {
				break
			}
		}
	}
	return posts
}
