package api

import (
	"bytes"
	"encoding/json"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/response"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestUserActions(t *testing.T) {
	files.Root = "../files/"
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	res := response.Info{}
	// No action
	got := userActions("", req, res)
	want := response.Info{}
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Can logoff
	authorized = false
	got = userActions("logon", req, res)
	want = response.Info{Msg: "Error"}
	if !strings.HasPrefix(got.Msg, want.Msg) {
		t.Errorf("Got %v want %v", got.Msg, want.Msg)
	}
	// Blocked from logoff when not authorized
	got = userActions("logoff", req, res)
	if got.Msg != "" {
		t.Errorf("Got %v want %v", got.Msg, "")
	}
	// Can logoff if authorized
	authorized = true
	got = userActions("logoff", req, res)
	if got.Msg != "ok" {
		t.Errorf("Got %v want %v", got.Msg, "ok")
	}
}

func TestPageActions(t *testing.T) {
	files.Root = "../files/"
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	res := response.Info{}
	// Save action blocked
	authorized = false
	got := pageActions("save", req, res)
	want := response.Info{}
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Save action authorized
	authorized = true
	info := page.EditInfo{ID: 777, Content: "xyz"}
	requestBody, _ := json.Marshal(info)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	got = pageActions("save", req, res)
	if got.Msg != "ok" {
		t.Errorf("Got %v want %v", got.Msg, "ok")
	}
	// Delete saved file
	page.Delete("test", 777)
	// Load file
	q := req.URL.Query()
	q.Add("pid", "1")
	req.URL.RawQuery = q.Encode()
	got = pageActions("load", req, res)
	want = response.Info{Msg: "ok"}
	if got.Msg != want.Msg {
		t.Errorf("Got %v want %v", got.Msg, want.Msg)
	}
	// Load file
	q.Set("pid", "9999")
	req.URL.RawQuery = q.Encode()
	got = pageActions("load", req, res)
	want = response.Info{Msg: "ok"}
	if got.Msg == want.Msg {
		t.Errorf("Got %v want %v", got.Msg, want.Msg)
	}
}

func TestPagesActions(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	res := response.Info{}
	// Actions blocked
	authorized = false
	got := pagesActions("save", req, res)
	want := response.Info{}
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Load data
	authorized = true
	got = pagesActions("load", req, res)
	if len(got.Data) < 1 {
		t.Errorf("Got no data trying to load")
	}
	requestBody := json.RawMessage(got.Data)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	got = pagesActions("save", req, res)
	got = pagesActions("load", req, res)
	b, _ := json.Marshal(requestBody)
	if got.Data != string(b) {
		t.Errorf("Got %v want %v", got.Data, string(b))
	}
}

func TestProcess(t *testing.T) {
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	_, _, response := Process(req, []string{}, map[string]string{})
	want := `{"ID":"","Data":"","Msg":""}`
	if response != want {
		t.Errorf("Got %s want %s", response, want)
	}
}
