package router

import (
	"gocms/pkg/files"
	"gocms/pkg/request"
	"testing"
)

func TestProcess(t *testing.T) {
	files.Root = "../files/"
	r := request.Info{Domain: "test", Route: "test_page"}
	Process(r)
}

func TestExtractSubRoutes(t *testing.T) {
	r := request.Info{Route: "/archive/a/b"}
	r, pageRoute := ExtractSubRoutes(r)
	want := "b"
	got := r.SubRoutes[1]
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
	got = pageRoute
	want = "/archive"
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}
