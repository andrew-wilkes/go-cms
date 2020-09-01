package router

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/user"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

// Request type
type Request struct {
	Domain   string
	Route    string
	Method   string
	GetArgs  map[string]string
	PostData map[string]json.RawMessage
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
		content = ReplaceTokens(r.Domain, template, page)
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
func ReplaceTokens(domain string, html string, page page.Info) string {
	year, month, day := time.Now().Date()
	if user.GetStatus().LoggedIn {
		html = strings.Replace(html, `#CSS#`, `<link rel="stylesheet" href="#HOST#/css/content-tools.min.css">`, 1)
		html = strings.Replace(html, `#SCRIPTS#`, getScripts(), 1)
	} else {
		html = strings.Replace(html, `#CSS#`, "", 1)
		html = strings.Replace(html, `#SCRIPTS#`, "", 1)
	}
	html = strings.ReplaceAll(html, `#HOST#`, domain)
	html = strings.ReplaceAll(html, `#TITLE#`, page.Title)
	html = strings.ReplaceAll(html, `#TOPMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#CONTENT#`, page.Content)
	html = strings.ReplaceAll(html, `#FOOTERMENU#`, page.Title)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#ARCHIVE#`, "Genetate archive content")
	html = strings.ReplaceAll(html, `#RECENT#`, "Genetate recent posts content")
	html = strings.ReplaceAll(html, `#CATEGORIES#`, "Genetate category list")
	return html
}

func getScripts() string {
	scripts := []string{"content-tools.min.js", "cloudinary.js", "editor.js", "axios.min.js", "common.js"}
	template := `<script src="#HOST#/%s"/>\n`
	html := ""
	for s := range scripts {
		html += fmt.Sprintf(template, s)
	}
	return html
}
