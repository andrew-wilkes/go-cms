package router

// This package builds the web page content according to the provided route and other request data

import (
	"gocms/pkg/content"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"strings"
)

// Process a request
func Process(r request.Info) ([]string, string) {
	r, pageRoute := ExtractSubRoutes(r)
	headers := []string{}
	html := ""
	page := page.GetByRoute(r.Domain, pageRoute, true)
	if page.ID == 0 {
		headers = append(headers, "HTTP/1.1 404 Not Found")
		html = "Page not found at: " + r.Route
	} else {
		if page.Template == "post" {
			// Avoid duplicate content issues with blog posts
			headers = append(headers, "rel: canonical")
		}
		template := files.GetTemplate(r.Domain, page.Template)
		html = content.ReplaceTokens(r, template, page)
	}
	return headers, html
}

// ExtractSubRoutes scans the route for special prefixes and uses the rest of the route to extract the subroutes
func ExtractSubRoutes(r request.Info) (request.Info, string) {
	stems := []string{"/archive"}
	pageRoute := r.Route
	for _, stem := range stems {
		if strings.HasPrefix(r.Route, stem) {
			tail := strings.Replace(r.Route, stem, "", 1)
			r.SubRoutes = strings.Split(tail, "/")[1:]
			pageRoute = stem
			break
		}
	}
	return r, pageRoute
}
