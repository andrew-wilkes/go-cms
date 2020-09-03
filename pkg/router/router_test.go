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
