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
	html = strings.ReplaceAll(html, `#TOPMENU#`, getPageLinks(p, baseLink))
	html = strings.ReplaceAll(html, `#CONTENT#`, p.Content)
	html = strings.ReplaceAll(html, `#FOOTERMENU#`, p.Title)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#ARCHIVE#`, "Generate archive content")
	html = strings.ReplaceAll(html, `#RECENT#`, "Generate recent posts content")
	html = strings.ReplaceAll(html, `#CATEGORIES#`, getCategoryLinks(p, baseLink))
	return html
}

func getScripts() string {
	scripts := []string{"content-tools.min.js", "cloudinary.js", "editor.js", "axios.min.js", "common.js"}
	template := `<script src="#HOST#/js/%s"/>\n`
	html := ""
	for s := range scripts {
		html += fmt.Sprintf(template, s)
	}
	return html
}

func getPageLinks(p page.Info, base string) string {
	// Return links to pages with a common parent
	links := ""
	pages := page.GetPages(p.Parent, page.Published)
	for _, item := range pages {
		links += fmt.Sprintf("<li>%s</li>\n", getHref(item, p.Route, base))
	}
	return links
}

func getCategoryLinks(p page.Info, base string) string {
	links := ""
	cats := page.GetCategoryPages(p.Category, page.Published)
	for _, c := range cats {
		if c.Parent > 0 {
			links += fmt.Sprintf("<li>%s</li>\n", getHref(c, p.Route, base))
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
