package content

import (
	"gocms/pkg/files"
	"gocms/pkg/page"
	"testing"
)

func Test(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	p := page.Get("test", 8, false)
	println(p.Title)
	println(GetBreadcrumbLinks("test", p, "base.com"))
}
