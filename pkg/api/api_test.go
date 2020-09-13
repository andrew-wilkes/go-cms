package api

import (
	"gocms/pkg/request"
	"testing"
)

func TestProcess(t *testing.T) {
	req := request.Info{}
	_, response := Process(req)
	want := `{"ID":"","Data":"","Msg":""}`
	if response != want {
		t.Errorf("Got %s want %s", response, want)
	}
}
