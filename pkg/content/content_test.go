package content

import (
	"crypto/md5"
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"testing"
)

func TestReplaceTokens(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	html := "#TITLE# #TOPMENU# #CONTENT# #FOOTERMENU# #DATE# #SCRIPTS#"
	page := page.Info{Title: "Test", Content: "Content"}
	content := ReplaceTokens(request.Info{Scheme: "http", Domain: "test"}, html, page)
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
