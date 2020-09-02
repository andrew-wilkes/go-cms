package content

import (
	"fmt"
	"gocms/pkg/page"
	"gocms/pkg/user"
	"strings"
	"time"
)

// ReplaceTokens in HTML
func ReplaceTokens(scheme string, domain string, html string, p page.Info) string {
	baseLink := fmt.Sprintf("%s://%s", scheme, domain)
	year, month, day := time.Now().Date()
	if user.GetStatus().LoggedIn {
		html = strings.Replace(html, `#CSS#`, `<link rel="stylesheet" href="#HOST#/css/content-tools.min.css">`, 1)
		html = strings.Replace(html, `#SCRIPTS#`, getScripts(), 1)
	} else {
		html = strings.Replace(html, `#CSS#`, "", 1)
		html = strings.Replace(html, `#SCRIPTS#`, "", 1)
	}
	html = strings.ReplaceAll(html, `#HOST#`, baseLink)
	html = strings.ReplaceAll(html, `#TITLE#`, p.Title)
	html = strings.ReplaceAll(html, `#TOPMENU#`, strings.Join(getPageLinks(p, baseLink, 2), ""))
	html = strings.ReplaceAll(html, `#CONTENT#`, p.Content)
	html = strings.ReplaceAll(html, `#FOOTERMENU#`, p.Title)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#ARCHIVE#`, "Generate archive content")
	html = strings.ReplaceAll(html, `#RECENT#`, "Generate recent posts content")
	html = strings.ReplaceAll(html, `#CATEGORIES#`, strings.Join(getCategoryLinks([]string{}, p, baseLink, 2), ""))
	return html
}

func getScripts() string {
	scripts := []string{"content-tools.min.js", "cloudinary.js", "editor.js", "axios.min.js", "common.js"}
	template := `<script src="#HOST#/js/%s"/>`
	html := []string{}
	for s := range scripts {
		html = append(html, fmt.Sprintf(template, s))
	}
	return strings.Join(html, "\n")
}

func getPageLinks(p page.Info, base string, depth int) []string {
	// Return links to pages with a common parent
	pages := page.GetPages(p.Parent, page.Published)
	return scanSubPages(pages, base, depth, p.Route)
}

func scanSubPages(pages []page.Info, base string, depth int, route string) []string {
	links := []string{}
	for _, item := range pages {
		subPageLinks := ""
		if depth > 1 {
			subPageLinks = strings.Join(getSubPageLinks(item, base, depth-1), "")
		}
		links = append(links, fmt.Sprintf("<li>%s%s</li>\n", getHref(item, route, base), subPageLinks))
	}
	return links
}

func getSubPageLinks(p page.Info, base string, depth int) []string {
	// Return links to pages with p.ID as parent
	pages := page.GetPages(p.ID, page.Published)
	return scanSubPages(pages, base, depth, "-")
}

func getCategoryLinks(links []string, p page.Info, base string, depth int) []string {
	cats := page.GetCategoryPages(p.ID, page.Published)
	for _, c := range cats {
		if c.Parent > 0 {
			// Subcategories are created by setting the Category value to the parent category page ID
			subCatLinks := ""
			if depth > 1 {
				subCatLinks = strings.Join(getCategoryLinks([]string{}, c, base, depth-1), "")
			}
			links = append(links, fmt.Sprintf("<li>%s%s</li>\n", getHref(c, p.Route, base), subCatLinks))
		}
	}
	return links
}

func getHref(p page.Info, route string, base string) string {
	href := p.Title
	if p.Route != route {
		href = fmt.Sprintf(`<a href="%s%s">%s</a>`, base, p.Route, href)
	}
	return href
}
