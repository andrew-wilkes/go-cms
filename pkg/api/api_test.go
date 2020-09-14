package api

import (
	"encoding/json"
	"gocms/pkg/files"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"strings"
	"testing"
	"time"
)

func TestSessionValid(t *testing.T) {
	// Session expired
	state = settings.Values{
		SessionExpiry: time.Now().AddDate(0, 0, -1),
		SessionID:     "aaa",
	}
	sessionValid("aaa")
	got := authorized
	want := false
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Session not expired, wrong id
	state.SessionExpiry = time.Now().AddDate(0, 0, 1)
	sessionValid("z")
	got = authorized
	want = false
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Session not expired, correct id
	sessionValid("aaa")
	got = authorized
	want = true
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
}

func TestUserActions(t *testing.T) {
	files.Root = "../files/"
	req := request.Info{Domain: "test"}
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
	req := request.Info{
		Domain:   "test",
		PostData: map[string]json.RawMessage{},
	}
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
	data, _ := json.Marshal(info)
	req.PostData["info"] = data
	got = pageActions("save", req, res)
	want = response.Info{Msg: "Error"}
	if got.Msg != "ok" {
		t.Errorf("Got %v want %v", got.Msg, "ok")
	}
	// Delete saved file
	page.Delete("test", 777)
	// Load file
	req.GetArgs = map[string]string{"pid": "1"}
	got = pageActions("load", req, res)
	want = response.Info{Msg: "ok"}
	if got.Msg != want.Msg {
		t.Errorf("Got %v want %v", got.Msg, want.Msg)
	}
	// Load file
	req.GetArgs = map[string]string{"pid": "9999"}
	got = pageActions("load", req, res)
	want = response.Info{Msg: "ok"}
	if got.Msg == want.Msg {
		t.Errorf("Got %v want %v", got.Msg, want.Msg)
	}
}

func TestPagesActions(t *testing.T) {
	files.Root = "../files/"
	page.LoadData("test")
	req := request.Info{Domain: "test"}
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
	rawData := json.RawMessage(got.Data)
	req.PostData = map[string]json.RawMessage{"pages": rawData}
	got = pagesActions("save", req, res)
	got = pagesActions("load", req, res)
	b, _ := json.Marshal(req.PostData["pages"])
	if got.Data != string(b) {
		t.Errorf("Got %v want %v", got.Data, string(b))
	}
}

func TestProcess(t *testing.T) {
	req := request.Info{}
	_, response := Process(req)
	want := `{"ID":"","Data":"","Msg":""}`
	if response != want {
		t.Errorf("Got %s want %s", response, want)
	}
}
