package router

import (
	"gocms/pkg/files"
	"testing"
)

func TestProcess(t *testing.T) {
	files.Root = "../files/"
	r := Request{"test", "test_page", make(map[string]string)}
	Process(r)
}

func TestGetTemplate(t *testing.T) {
	files.Root = "../files/"
	GetTemplate("test", "home")
}
