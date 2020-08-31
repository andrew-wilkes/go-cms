package router

import (
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

// Request type
type Request struct {
	Domain string
	Route  string
	Params map[string]string
}

// Process a request
func Process(r Request) ([]string, string) {
	headers := []string{}
	content := ""
	page := page.GetByRoute(r.Domain, r.Route, true)
	if page.ID == 0 {
		headers = append(headers, "HTTP/1.1 404 Not Found")
		content = "Page not found at: " + r.Route
	} else {
		template := GetTemplate(r.Domain, page.Template)
		content = ReplaceTokens(template, page)
	}
	return headers, content
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
	year, month, day := time.Now().Date()
	html = strings.ReplaceAll(html, `#TITLE#`, page.Title)
	html = strings.ReplaceAll(html, `#TOPMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#CONTENT#`, page.Content)
	html = strings.ReplaceAll(html, `#FOOTERMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#SCRIPTS#`, page.Title)
	return html
}
