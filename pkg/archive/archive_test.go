package archive

import (
	"fmt"
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
	gotr := years[2000]
	wantr := Record{Count: 20}
	if gotr.Count != wantr.Count {
		t.Errorf("Got %d want %d", gotr.Count, wantr.Count)
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
	gotr := months[int(time.February)]
	wantr := Record{Count: 5}
	if gotr.Count != wantr.Count {
		t.Errorf("Got %d want %d", gotr.Count, wantr.Count)
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

func TestGetKeysInOrder(t *testing.T) {
	// Input map
	m := make(map[int]Record)
	m[1] = Record{Count: 11}
	m[2] = Record{Count: 12}
	m[0] = Record{Count: 10}
	sortedKeys := GetKeysInOrder(m)
	got := fmt.Sprintf("%d%d%d", sortedKeys[0], sortedKeys[1], sortedKeys[2])
	want := "012"
	if got != want {
		t.Errorf("Got %s want %s", got, want)
	}
}
