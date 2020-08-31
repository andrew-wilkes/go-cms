package router

import (
	"gocms/pkg/files"
	"gocms/pkg/page"
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

func TestReplaceTokens(t *testing.T) {
	html := "#TITLE# #TOPMENU# #CONTENT# #FOOTERMENU# #DATE# #SCRIPTS#"
	page := page.Info{Title: "Test", Content: "Content"}
	if ReplaceTokens(html, page) != "Test Test Content Test Test Test" {
		t.Fail()
	}
}
