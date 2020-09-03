package router

import (
	"gocms/pkg/files"
	"testing"
)

func TestProcess(t *testing.T) {
	files.Root = "../files/"
	r := Request{Domain: "test", Route: "test_page"}
	Process(r)
}

func TestGetTemplate(t *testing.T) {
	files.Root = "../files/"
	GetTemplate("test", "home")
}

func TestExtractSubRoutes(t *testing.T) {
	r := Request{Route: "/archive/a/b"}
	r = ExtractSubRoutes(r)
	want := "b"
	got := r.SubRoutes[1]
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}
