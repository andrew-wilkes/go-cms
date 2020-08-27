package category

import (
	"gocms/pkg/page"
)

// Info type
type Info struct {
	ID     int
	Parent int
	Depth  int
	Title  string
	Route  string
	Count  int
}

// GetPosts - get a list of posts within a category
func GetPosts(id int) []page.Info {
	return []page.Info{
		page.Get(1, false),
	}
}

// List - Get a list of all categories
func List(id int) []Info {
	return []Info{
		Info{ID: 1, Title: "page title", Route: "/page-route", Count: 12},
	}
}
