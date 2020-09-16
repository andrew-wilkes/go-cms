package router

import (
	"bytes"
	"gocms/pkg/files"
	"net/http"
	"testing"
)

func TestProcess(t *testing.T) {
	files.Root = "../files/"
	req, _ := http.NewRequest("GET", "", bytes.NewBufferString(""))
	req.Host = "test"
	req.RequestURI = "/test_page"
	Process(req)
}

func TestExtractSubRoutes(t *testing.T) {
	req, _ := http.NewRequest("GET", "", bytes.NewBufferString(""))
	req.Host = "base.com"
	req.RequestURI = "/archive/a/b"
	subRoutes, pageRoute := ExtractSubRoutes(req)
	want := "b"
	got := subRoutes[1]
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
	got = pageRoute
	want = "/archive"
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}
