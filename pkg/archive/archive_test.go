package archive

import (
	"gocms/pkg/files"
	"gocms/pkg/page"
	"testing"
	"time"
)

func TestGetYears(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	years := GetYears()
	got := len(years)
	want := 5
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

func TestGetMonths(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	months := GetMonths(2000)
	got := len(months)
	want := 4
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

func TestGetDays(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	months := GetDays(2000, time.February)
	got := len(months)
	want := 4
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}
