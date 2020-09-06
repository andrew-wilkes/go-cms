package content

import (
	"fmt"
	"gocms/pkg/archive"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"gocms/pkg/user"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReplaceTokens in HTML template
func ReplaceTokens(r request.Info, html string, p page.Info) string {
	baseURL := fmt.Sprintf("%s://%s", r.Scheme, r.Domain)
	year, month, day := time.Now().Date()
	if user.GetStatus().LoggedIn {
		html = strings.Replace(html, `#CSS#`, `<link rel="stylesheet" href="#HOST#/css/content-tools.min.css">`, 1)
		html = strings.Replace(html, `#SCRIPTS#`, getScripts(), 1)
	} else {
		html = strings.Replace(html, `#CSS#`, "", 1)
		html = strings.Replace(html, `#SCRIPTS#`, "", 1)
	}
	html = strings.ReplaceAll(html, `#HOST#`, baseURL)
	html = strings.ReplaceAll(html, `#HOME#`, getHref(page.GetPages(0, page.Published)[0], p.Route, baseURL))
	html = strings.ReplaceAll(html, `#BREADCRUMB#`, GetBreadcrumbLinks(r.Domain, p, baseURL))
	html = strings.ReplaceAll(html, `#CONTENT#`, p.Content)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#ARCHIVE#`, generateArchive(r, baseURL))
	html = strings.ReplaceAll(html, `#RECENT#`, "Generate recent posts content")
	html = strings.ReplaceAll(html, `#PAGESINCATEGORY#`, getPagesInCategory(p, baseURL))
	html = strings.ReplaceAll(html, `#TITLE#`, p.Title)
	html = addMenus(html, baseURL)
	html = addPageLinks(html, baseURL, p)
	html = addCategoryLinks(html, baseURL, p)
	html = addRecentPostsLinks(html, baseURL)
	return html
}

func addMenus(html string, baseURL string) string {
	re, _ := regexp.Compile(`#(\w+)_MENU#`) // e.g. SIDE_MENU
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		for _, token := range m {
			menu := makeHTMLList(page.GetPagesInMenu(strings.ToLower(token[1])), baseURL)
			html = strings.ReplaceAll(html, token[0], menu)
		}
	}
	return html
}

func addPageLinks(html string, baseURL string, currentPage page.Info) string {
	re, _ := regexp.Compile(`#PAGES_(\d)#`) // e.g. #PAGES_2# to get a list that is 2 levels deep
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		for _, token := range m {
			list := strings.Join(getPageLinks(currentPage, baseURL, getInt(token[1])), "\n")
			html = strings.ReplaceAll(html, token[0], list)
		}
	}
	return html
}

func addCategoryLinks(html string, baseURL string, currentPage page.Info) string {
	re, _ := regexp.Compile(`#CATEGORIES_(\d)#`) // e.g. #CATEGORIES_2# to get a list that is 2 levels deep
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		for _, token := range m {
			list := strings.Join(getCategoryLinks(currentPage, baseURL, getInt(token[1])), "\n")
			html = strings.ReplaceAll(html, token[0], list)
		}
	}
	return html
}

func addRecentPostsLinks(html string, baseURL string) string {
	re, _ := regexp.Compile(`#RECENT_(\d+)#`) // e.g. #RECENT_2# to get a list that is 2 levels deep
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		for _, token := range m {
			list := makeHTMLList(page.GetRecentPosts(getInt(token[1])), baseURL)
			html = strings.ReplaceAll(html, token[0], list)
		}
	}
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

func generateArchive(r request.Info, baseURL string) string {
	// The outer HTML will likely be UL tags in the template since a css class may be applied to it
	var links []string
	switch len(r.SubRoutes) {
	case 0:
		links = getYearArchiveLinks(baseURL)
	case 1:
		links = getMonthArchiveLinks(getInt(r.SubRoutes[0]), baseURL)
	case 2:
		links = getDayArchiveLinks(getInt(r.SubRoutes[0]), time.Month(getInt(r.SubRoutes[1])), baseURL)
	}
	return strings.Join(links, "\n")
}

func getInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i // Return 0 if there is an error
}

func getYearArchiveLinks(baseURL string) []string {
	links := []string{}
	for year, count := range archive.GetYears() {
		links = append(links, fmt.Sprintf(`<li>%d (%d)</li>`, year, count))
		links = append(links, "<ul>")
		links = append(links, getMonthArchiveLinks(year, baseURL)...)
		links = append(links, "</ul>")
	}
	return links
}

func getMonthArchiveLinks(year int, baseURL string) []string {
	links := []string{}
	for month, count := range archive.GetMonths(year) {
		links = append(links, fmt.Sprintf(`<li>%s (%d)</li>`, month, count))
		links = append(links, "<ul>")
		links = append(links, getDayArchiveLinks(year, month, baseURL)...)
		links = append(links, "</ul>")
	}
	return links
}

func getDayArchiveLinks(year int, month time.Month, baseURL string) []string {
	links := []string{}
	for _, posts := range archive.GetDays(year, month) {
		for _, p := range posts {
			links = append(links, fmt.Sprintf("<li>%s %s</li>", p.PubDate.Format("2006-01-02 15:04 Monday"), getHref(p, "-", baseURL)))
		}
	}
	return links
}

func getPostLinks(posts []page.Info) []string {
	links := []string{}
	return links
}

func getPageLinks(p page.Info, baseURL string, depth int) []string {
	// Return links to pages with a common parent
	pages := page.GetPages(p.Parent, page.Published)
	return scanSubPages(pages, baseURL, depth, p.Route)
}

func scanSubPages(pages []page.Info, baseURL string, depth int, route string) []string {
	links := []string{}
	for _, item := range pages {
		subPageLinks := ""
		if depth > 1 {
			subPageLinks = fmt.Sprintf("\n<ul>\n%s</ul>\n", strings.Join(getSubPageLinks(item, baseURL, depth-1), ""))
		}
		links = append(links, fmt.Sprintf("<li>%s%s</li>\n", getHref(item, route, baseURL), subPageLinks))
	}
	return links
}

func getSubPageLinks(p page.Info, base string, depth int) []string {
	// Return links to pages with p.ID as parent
	pages := page.GetPages(p.ID, page.Published)
	return scanSubPages(pages, base, depth, "-")
}

func getCategoryLinks(p page.Info, base string, depth int) []string {
	links := []string{}
	cats := page.GetCategoryPages(p.ID, page.Published)
	for _, c := range cats {
		if c.Parent > 0 {
			// Subcategories are created by setting the Category value to the parent category page ID
			subCatLinks := ""
			if depth > 1 {
				subCatLinks = fmt.Sprintf("\n<ul>%s</ul>\n", strings.Join(getCategoryLinks(c, base, depth-1), ""))
			}
			links = append(links, fmt.Sprintf("<li>%s%s</li>\n", getHref(c, p.Route, base), subCatLinks))
		}
	}
	return links
}

func getPagesInCategory(p page.Info, baseURL string) string {
	return makeHTMLList(page.GetPagesInCategory(p.ID), baseURL)
}

func makeHTMLList(pages []page.Info, baseURL string) string {
	links := []string{}
	for _, item := range pages {
		links = append(links, fmt.Sprintf("<li>%s</li>\n", getHref(item, "-", baseURL)))
	}
	return strings.Join(links, "\n")
}

func getHref(p page.Info, route string, baseURL string) string {
	href := p.Title
	if p.Route != route {
		href = fmt.Sprintf(`<a href="%s%s">%s</a>`, baseURL, p.Route, href)
	}
	return href
}

// GetBreadcrumbLinks returns a trail of links from the home page to the current page
func GetBreadcrumbLinks(domain string, p page.Info, baseURL string) string {
	pages := []page.Info{p}
	pid := p.Parent
	for pid > 0 {
		parentPage := page.GetByID(domain, pid, false)
		pages = append([]page.Info{parentPage}, pages...)
		pid = parentPage.Parent
	}
	links := []string{}
	for _, bp := range pages {
		links = append(links, getHref(bp, p.Route, baseURL))
	}
	return strings.Join(links, " > ")
}
