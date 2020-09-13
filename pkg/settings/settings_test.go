package settings

import (
	"gocms/pkg/files"
	"testing"
)

func TestGet(t *testing.T) {
	files.Root = "../files/"
	got := Get("test").loaded
	want := true
	if got != want {
		t.Errorf("Want %v got %v", want, got)
	}
	got = Get("test").loaded
}

func TestSet(t *testing.T) {
	files.Root = "../files/"
	Set(Values{UserName: "Admin"}, "test")
	got := values.UserName
	want := "Admin"
	if got != want {
		t.Errorf("Want %v got %v", want, got)
	}
}
