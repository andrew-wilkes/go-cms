package content

import (
	"fmt"
	"gocms/pkg/page"
	"gocms/pkg/user"
	"strings"
	"time"
)

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
	html = strings.ReplaceAll(html, `#ARCHIVE#`, "Generate archive content")
	html = strings.ReplaceAll(html, `#RECENT#`, "Generate recent posts content")
	html = strings.ReplaceAll(html, `#CATEGORIES#`, "Generate category list")
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
