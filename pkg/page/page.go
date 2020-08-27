package page

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
	Timestamp   int
}

// Get a page or post
func Get(id int, getContent bool) Info {
	return Info{Title: "page title", Route: "/page-route", Timestamp: 3456}
}

// Save a page or post
func Save(info Info, saveContent bool) int {
	return 1
}

// List - Get a list of all pages
func List() []Info {
	return []Info{Info{}, Info{}}
}
