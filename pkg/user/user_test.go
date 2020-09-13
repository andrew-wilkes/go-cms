package user

import (
	"encoding/json"
	"gocms/pkg/files"
	"gocms/pkg/request"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"strings"
	"testing"
)

func TestLogOn(t *testing.T) {
	files.Root = "../files/"
	resp := response.Info{}
	// Test with bad data
	req := request.Info{
		Domain:   "test",
		PostData: map[string]json.RawMessage{},
	}
	resp = LogOn(req, resp)
	got := resp.ID
	want := ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	got = resp.Msg
	want = "Error decoding data!"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	settings.Set(settings.Values{Email: "x@y.com"}, "test")
	// Test with unknown user (email)
	resp = response.Info{}
	testUser := Credentials{Name: "Admin", Email: "a@b.com", Pass: "pwd"}
	data, _ := json.Marshal(testUser)
	req.PostData["user"] = data
	resp = LogOn(req, resp)
	got = resp.ID
	want = ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	got = resp.Msg
	want = "Unknown user!"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	// Test with known user but wrong password
	resp = response.Info{}
	settings.Set(settings.Values{UserName: "Admin", Email: "a@b.com", Password: hash("abc")}, "test")
	resp = LogOn(req, resp)
	got = resp.ID
	want = ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	got = resp.Msg
	want = "Wrong password!"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	// Test with known user and correct password
	resp = response.Info{}
	settings.Set(settings.Values{UserName: "Admin", Email: "a@b.com", Password: hash("pwd")}, "test")
	resp = LogOn(req, resp)
	goti := len(resp.ID)
	wanti := 20
	if wanti != goti {
		t.Errorf("Got %d want %d", goti, wanti)
	}
	got = resp.Msg
	want = "ok"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
}

func TestLogOff(t *testing.T) {
	files.Root = "../files/"
	resp := response.Info{}
	req := request.Info{
		Domain: "test",
	}
	resp = LogOff(req, resp)
	got := resp.Msg
	want := "ok"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	info := settings.Get(req.Domain)
	got = info.SessionID
	want = ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
}

func TestRegister(t *testing.T) {
	files.Root = "../files/"
	resp := response.Info{}
	// Test with bad data
	req := request.Info{
		Domain:   "test",
		PostData: map[string]json.RawMessage{},
	}
	resp = Register(req, resp)
	got := resp.ID
	want := ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	got = resp.Msg
	want = "Error decoding data!"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	// Test with short password
	resp = response.Info{}
	testUser := Credentials{Name: "Admin", Email: "a@b.com", Pass: "pwd"}
	data, _ := json.Marshal(testUser)
	req.PostData["user"] = data
	resp = Register(req, resp)
	got = resp.ID
	want = ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	if !strings.HasPrefix(resp.Msg, "Please") {
		t.Errorf("Got %s", resp.Msg)
	}
	testUser = Credentials{Name: "Admin", Email: "a@b.com", Pass: "123456"}
	data, _ = json.Marshal(testUser)
	req.PostData["user"] = data
	resp = Register(req, resp)
	got = resp.Msg
	want = "ok"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
}
