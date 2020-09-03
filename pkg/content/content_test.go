package content

import (
	"gocms/pkg/files"
	"gocms/pkg/page"
	"testing"
)

func TestGetBreadcrumbLinks(t *testing.T) {
	files.Root = "../files/"
	p := page.GetByID("test", 8, false)
	want := `<a href="base.com/">Page 0</a> > <a href="base.com/0-1">Page 0-1</a> > Page 0-1-0`
	got := GetBreadcrumbLinks("test", p, "base.com")
	if got != want {
		t.Errorf("Want %s got %s", want, got)
	}
}
