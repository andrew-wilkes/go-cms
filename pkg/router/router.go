package router

// This package builds the web page content according to the provided route and other request data

import (
	"encoding/json"
	"gocms/pkg/content"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"io/ioutil"
	"regexp"
	"strings"
)

// Request type
type Request struct {
	Domain   string
	Route    string
	Method   string
	Scheme   string
	GetArgs  map[string]string
	PostData map[string]json.RawMessage
}

// Process a request
func Process(r Request) ([]string, string) {
	headers := []string{}
	html := ""
	page := page.GetByRoute(r.Domain, r.Route, true)
	if page.ID == 0 {
		headers = append(headers, "HTTP/1.1 404 Not Found")
		html = "Page not found at: " + r.Route
	} else {
		template := GetTemplate(r.Domain, page.Template)
		html = content.ReplaceTokens(r.Scheme, r.Domain, template, page)
	}
	return headers, html
}

// GetTemplate for page (this is a recursive function to allow for nesting of templates)
func GetTemplate(domain string, templateName string) string {
	data, err := ioutil.ReadFile(files.Root + domain + "/templates/" + templateName + ".html")
	if err != nil {
		panic(err)
	}
	template := string(data)
	// Look for sub-template tokens
	re, _ := regexp.Compile(`#T_(\w+)#`) // An example token targeting the footer.html template is: #T_FOOTER#
	m := re.FindAllStringSubmatch(template, -1)
	if m != nil {
		for _, sub := range m {
			template = strings.ReplaceAll(template, sub[0], GetTemplate(domain, strings.ToLower(sub[1])))
		}
	}
	return template
}
