package router

import (
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"io/ioutil"
	"regexp"
	"strings"
)

// Request type
type Request struct {
	Domain string
	Route  string
	Params map[string]string
}

// Process a request
func Process(r Request) {
	page := page.GetByRoute(r.Domain, r.Route, true)
	if page.ID == 0 {
		fmt.Println("Page not found!")
	} else {
		template := GetTemplate(r.Domain, page.Template)
		fmt.Println(template)
	}
}

// GetTemplate for page
func GetTemplate(domain string, template string) string {
	data, err := ioutil.ReadFile(files.Root + domain + "/templates/" + template + ".html")
	if err != nil {
		panic(err)
	}
	content := string(data)
	// Look for sub-template tags
	re, _ := regexp.Compile(`#T_(.+)#`)
	m := re.FindAllStringSubmatch(content, -1)
	if m != nil {
		for _, sub := range m {
			content = strings.ReplaceAll(content, sub[0], GetTemplate(domain, strings.ToLower(sub[1])))
		}
	}
	return content
}

// ReplaceTokens in HTML
func ReplaceTokens(html string, page page.Info) string {
	html = strings.ReplaceAll(html, `#TITLE#`, page.Title)
	html = strings.ReplaceAll(html, `#TOPMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#CONTENT#`, page.Content)
	html = strings.ReplaceAll(html, `#FOOTERMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#DATE#`, page.Title)
	html = strings.ReplaceAll(html, `#SCRIPTS#`, page.Title)
	return html
}
