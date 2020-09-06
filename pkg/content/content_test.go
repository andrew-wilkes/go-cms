package content

import (
	"crypto/md5"
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"strings"
	"testing"
)

func TestReplaceTokens(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	html := "#TITLE# #TOPMENU# #CONTENT# #FOOTERMENU# #DATE# #SCRIPTS#"
	page := page.Info{Title: "Test", Content: "Content"}
	content := ReplaceTokens(request.Info{Scheme: "http", Domain: "test"}, html, page)
	println(content)
	got := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	want := "ccabdae6e8a3dce76ac249c82bb0ed26"
	if got != want {
		t.Errorf("Want %s got %s", want, got)
	}
}

func TestGetBreadcrumbLinks(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	p := page.GetByID("test", 8, false)
	want := `<a href="base.com/">Page 0</a> > <a href="base.com/0-1">Page 0-1</a> > Page 0-1-0`
	got := GetBreadcrumbLinks("test", p, "base.com")
	if got != want {
		t.Errorf("Want %s got %s", want, got)
	}
}

func TestGetDayArchiveLinks(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	var got string
	for _, day := range getDayArchiveLinks(2015, 2, "test.com") {
		println(day)
		got = day
	}
	want := `<li>2015`
	if !strings.HasPrefix(got, want) {
		t.Errorf("Want %s got %s", want, got)
	}
}

func TestGetMonthArchiveLinks(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	var got int
	for _, month := range getMonthArchiveLinks(2015, "test.com") {
		println(month)
		got++
	}
	want := 36
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGetYearArchiveLinks(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	var got string
	for _, month := range getYearArchiveLinks("test.com") {
		println(month)
		got = month
	}
	want := `</ul>`
	if !strings.HasPrefix(got, want) {
		t.Errorf("Want %s got %s", want, got)
	}
}

func TestGetInt(t *testing.T) {
	want := 0
	got := getInt("a")
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
	want = 2015
	got = getInt("2015")
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGenerateArchive(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	r := request.Info{}
	baseURL := "test.com"
	got := len(generateArchive(r, baseURL))
	want := 6882
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
	r.SubRoutes = append(r.SubRoutes, "2015")
	got = len(generateArchive(r, baseURL))
	want = 1682
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
	r.SubRoutes = append(r.SubRoutes, "01")
	got = len(generateArchive(r, baseURL))
	want = 65
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}

func TestGetPagesInCategory(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	pages := getPagesInCategory(page.Info{ID: 1}, "test.com")
	got := len(pages)
	want := 3199
	if got != want {
		t.Errorf("Want %d got %d", want, got)
	}
}
