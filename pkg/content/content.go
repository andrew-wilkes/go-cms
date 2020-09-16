package content

import (
	"fmt"
	"gocms/pkg/archive"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/user"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReplaceTokens in HTML template
func ReplaceTokens(req *http.Request, html string, p page.Info, subRoutes []string) string {
	baseURL := fmt.Sprintf("%s://%s", req.URL.Scheme, req.Host)
	year, month, day := time.Now().Date()
	if user.SessionValid(req.URL.Query()["id"], req.Host) {
		html = strings.Replace(html, `#CSS#`, `<link rel="stylesheet" href="#HOST#/css/content-tools.min.css">`, 1)
		html = strings.Replace(html, `#SCRIPTS#`, getScripts(), 1)
	} else {
		html = strings.Replace(html, `#CSS#`, "", 1)
		html = strings.Replace(html, `#SCRIPTS#`, "", 1)
	}
	html = strings.ReplaceAll(html, `#HOST#`, baseURL)
	html = strings.ReplaceAll(html, `#ID#`, fmt.Sprintf("%d", p.ID))
	// The home page is associated with an ID of 1
	html = strings.ReplaceAll(html, `#HOME#`, getHref(page.GetByID(req.Host, 1, false), p.Route, baseURL))
	html = strings.ReplaceAll(html, `#BREADCRUMB#`, GetBreadcrumbLinks(req, p, baseURL))
	html = strings.ReplaceAll(html, `#CONTENT#`, p.Content)
	html = strings.ReplaceAll(html, `#DAY#`, fmt.Sprint(day))
	html = strings.ReplaceAll(html, `#MONTH#`, fmt.Sprint(month))
	html = strings.ReplaceAll(html, `#YEAR#`, fmt.Sprint(year))
	html = strings.ReplaceAll(html, `#ARCHIVE#`, generateArchive(req, baseURL, &p, subRoutes))
	html = strings.ReplaceAll(html, `#PAGESINCATEGORY#`, getPagesInCategory(p, baseURL))
	html = strings.ReplaceAll(html, `#TITLE#`, p.Title)
	html = strings.ReplaceAll(html, "#DESCRIPTION#", p.Description)
	html = strings.Replace(html, "#COMMENTS#", "", 1)
	html = addMenus(html, baseURL, req.RequestURI)
	html = addPageLinks(html, baseURL, p)
	html = addCategoryLinks(html, baseURL, p)
	html = addRecentPostsLinks(html, baseURL)
	html = addPosts(req.Host, html, baseURL)
	html = strings.Replace(html, "#MORE#", "", 1)
	return html
}

func addMenus(html string, baseURL string, route string) string {
	re, _ := regexp.Compile(`#(\w+)_MENU#`) // e.g. SIDE_MENU
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		m2 := ""
		for _, token := range m {
			m1 := strings.ToLower(token[1])
			// Avoid repetition
			if m1 != m2 {
				menu := makeHTMLList(page.GetPagesInMenu(m1), route, baseURL)
				html = strings.ReplaceAll(html, token[0], menu)
			}
			m2 = m1
		}
	}
	return html
}

func addPageLinks(html string, baseURL string, currentPage page.Info) string {
	re, _ := regexp.Compile(`#PAGES_(\d)#`) // e.g. #PAGES_2# to get a list that is 2 levels deep
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		lastDepth := 0
		for _, token := range m {
			depth := getInt(token[1])
			// Avoid repetition
			if lastDepth != depth {
				list := strings.Join(getPageLinks(currentPage, baseURL, getInt(token[1]), "page"), "\n")
				html = strings.ReplaceAll(html, token[0], list)
			}
			lastDepth = depth
		}
	}
	return html
}

func addCategoryLinks(html string, baseURL string, currentPage page.Info) string {
	re, _ := regexp.Compile(`#CATEGORIES_(\d)#`) // e.g. #CATEGORIES_2# to get a list that is 2 levels deep
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		lastDepth := 0
		for _, token := range m {
			depth := getInt(token[1])
			// Avoid repetition
			if lastDepth != depth {
				list := strings.Join(getCategoryLinks(currentPage, baseURL, depth), "\n")
				html = strings.ReplaceAll(html, token[0], list)
			}
			lastDepth = depth
		}
	}
	return html
}

func addRecentPostsLinks(html string, baseURL string) string {
	re, _ := regexp.Compile(`#RECENT_(\d+)#`) // e.g. #RECENT_10# to get a list of 10 posts
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		lastDepth := 0
		for _, token := range m {
			depth := getInt(token[1])
			// Avoid repetition
			if lastDepth != depth {
				list := makeHTMLList(page.GetRecentPosts(getInt(token[1]), false, ""), "-", baseURL)
				html = strings.ReplaceAll(html, token[0], list)
			}
			lastDepth = depth
		}
	}
	return html
}

func addPosts(domain string, html string, baseURL string) string {
	re, _ := regexp.Compile(`#POSTS_(\d+)#`) // e.g. #POSTS_1# to insert 1 post
	m := re.FindAllStringSubmatch(html, -1)
	if m != nil {
		token := m[0]
		template := files.GetTemplate(domain, "inline-post")
		posts := page.GetRecentPosts(getInt(token[1]), true, domain)
		content := []string{}
		for _, p := range posts {
			year, month, day := p.PubDate.Date()
			t := template
			t = strings.Replace(t, "#YEAR#", fmt.Sprint(year), 1)
			t = strings.Replace(t, "#MONTH#", fmt.Sprint(month), 1)
			t = strings.Replace(t, "#DAY#", fmt.Sprint(day), 1)
			t = strings.Replace(t, "#TITLE#", p.Title, 1)
			t = strings.Replace(t, "#DESCRIPTION#", p.Description, 1)
			subStr := strings.Split(p.Content, "#MORE#")
			if len(subStr) > 1 {
				p.Title = "Read more ..."
				p.Content = fmt.Sprintf(`%s<p class="more">%s</p>`, subStr[0], getHref(p, "", baseURL))
			}
			t = strings.Replace(t, "#CONTENT#", p.Content, 1)
			content = append(content, t)
		}
		html = strings.Replace(html, token[0], strings.Join(content, "\n"), 1)
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

func generateArchive(req *http.Request, baseURL string, p *page.Info, subRoutes []string) string {
	// The outer HTML will likely be UL tags in the template since a css class may be applied to it
	var links []string
	switch len(subRoutes) {
	case 0:
		links = getYearArchiveLinks(baseURL)
	case 1:
		year := getInt(subRoutes[0])
		p.Title = fmt.Sprintf("Archives for %d", year)
		links = getMonthArchiveLinks(year, baseURL)
	case 2:
		year := getInt(subRoutes[0])
		month := time.Month(getInt(subRoutes[1]))
		p.Title = fmt.Sprintf("Archives for %s %d", month, year)
		links = getDayArchiveLinks(year, month, baseURL)
	}
	return strings.Join(links, "\n")
}

func getInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i // Return 0 if there is an error
}

func getYearArchiveLinks(baseURL string) []string {
	links := []string{}
	years := archive.GetYears()
	keys := archive.GetKeysInOrder(years)
	// Loop from most recent year
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]
		year := years[key]
		links = append(links, fmt.Sprintf(`<li><a href="%s/archive/%d">%d</a> (%d)</li>`, baseURL, key, key, year.Count))
		links = append(links, "<ul>")
		links = append(links, getMonthArchiveLinks(key, baseURL)...)
		links = append(links, "</ul>")
	}
	return links
}

func getMonthArchiveLinks(year int, baseURL string) []string {
	links := []string{}
	months := archive.GetMonths(year)
	keys := archive.GetKeysInOrder(months)
	for _, key := range keys {
		month := months[key]
		links = append(links, fmt.Sprintf(`<li><a href="%s/archive/%d/%d">%s</a> (%d)</li>`, baseURL, year, key, month.Month, month.Count))
		links = append(links, "<ul>")
		links = append(links, getDayArchiveLinks(year, month.Month, baseURL)...)
		links = append(links, "</ul>")
	}
	return links
}

func getDayArchiveLinks(year int, month time.Month, baseURL string) []string {
	links := []string{}
	days := archive.GetDays(year, month)
	keys := archive.GetKeysInOrder(days)
	for _, key := range keys {
		for _, p := range days[key].Posts {
			links = append(links, fmt.Sprintf("<li>%s %s</li>", getHref(p, "-", baseURL), p.PubDate.Format("Monday _2 15:04")))
		}
	}
	return links
}

func getPageLinks(p page.Info, baseURL string, depth int, template string) []string {
	// Return links to children of home page, else return links to sibling pages
	id := 1 // ID of home page
	if p.ID != id {
		id = p.Parent
	}
	pages := page.GetPages(id, page.Published, template)
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
	pages := page.GetPages(p.ID, page.Published, p.Template)
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
			links = append(links, fmt.Sprintf("<li>%s  (%d)%s</li>\n", getHref(c, p.Route, base), page.CountPagesInCategory(c.ID), subCatLinks))
		}
	}
	return links
}

func getPagesInCategory(p page.Info, baseURL string) string {
	return makeHTMLList(page.GetPagesInCategory(p.ID), "-", baseURL)
}

func makeHTMLList(pages []page.Info, route string, baseURL string) string {
	links := []string{}
	for _, item := range pages {
		links = append(links, fmt.Sprintf("<li>%s</li>\n", getHref(item, route, baseURL)))
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
func GetBreadcrumbLinks(req *http.Request, p page.Info, baseURL string) string {
	pages := []page.Info{p}
	pid := p.Parent
	for pid > 0 {
		parentPage := page.GetByID(req.Host, pid, false)
		pages = append([]page.Info{parentPage}, pages...)
		pid = parentPage.Parent
	}
	links := []string{}
	for _, bp := range pages {
		links = append(links, getHref(bp, req.RequestURI, baseURL))
	}
	return strings.Join(links, " > ")
}
