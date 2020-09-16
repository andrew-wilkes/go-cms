package router

// This package builds the web page content according to the provided route and other request data

import (
	"gocms/pkg/api"
	"gocms/pkg/content"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"net/http"
	"strings"
)

// Process a request
func Process(req *http.Request) (int, map[string]string, string) {
	subRoutes, pageRoute := ExtractSubRoutes(req)
	status := http.StatusOK
	headers := map[string]string{}
	domain := req.Host
	html := ""
	if pageRoute == "/api" {
		return api.Process(req, subRoutes)
	}
	page := page.GetByRoute(domain, pageRoute, true)
	if page.ID == 0 {
		status = http.StatusNotFound
		html = "Page not found at: " + pageRoute
	} else {
		if page.Template == "post" {
			// Avoid duplicate content issues with blog posts
			headers["rel"] = "canonical"
		}
		template := files.GetTemplate(domain, page.Template)
		html = content.ReplaceTokens(req, template, page, subRoutes)
	}
	return status, headers, html
}

// ExtractSubRoutes scans the route for special prefixes and uses the rest of the route to extract the subroutes
func ExtractSubRoutes(req *http.Request) ([]string, string) {
	subRoutes := []string{}
	stems := []string{"/archive", "/api"}
	pageRoute := req.RequestURI
	for _, stem := range stems {
		if strings.HasPrefix(pageRoute, stem) {
			tail := strings.Replace(pageRoute, stem, "", 1)
			subRoutes = strings.Split(tail, "/")[1:]
			pageRoute = stem
			break
		}
	}
	return subRoutes, pageRoute
}
