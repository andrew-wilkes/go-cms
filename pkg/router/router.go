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
	Domain    string
	Route     string
	SubRoutes []string
	Method    string
	Scheme    string
	GetArgs   map[string]string
	PostData  map[string]json.RawMessage
}

// Process a request
func Process(r Request) ([]string, string) {
	r = ExtractSubRoutes(r)
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

// ExtractSubRoutes scans the route for special prefixes and uses the rest of the route to extract the subroutes
func ExtractSubRoutes(r Request) Request {
	stems := []string{"/archive"}
	for _, stem := range stems {
		if strings.HasPrefix(r.Route, stem) {
			tail := strings.Replace(r.Route, stem, "", 1)
			r.SubRoutes = strings.Split(tail, "/")[1:]
			r.Route = stem
			break
		}
	}
	return r
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
