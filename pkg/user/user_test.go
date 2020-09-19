package user

import (
	"bytes"
	"encoding/json"
	"gocms/pkg/files"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestLogOn(t *testing.T) {
	files.Root = "../files/"
	resp := response.Info{}
	// Test with bad data
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	if err != nil {
		t.Errorf("%s", err)
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
	requestBody, err := json.Marshal(testUser)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
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
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
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
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
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
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	if err != nil {
		t.Errorf("%s", err)
	}
	resp = LogOff(req, resp)
	got := resp.Msg
	want := "ok"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	info := settings.Get(req.Host)
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
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Host = "test"
	if err != nil {
		t.Errorf("%s", err)
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
	requestBody, err := json.Marshal(testUser)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	resp = Register(req, resp)
	got = resp.ID
	want = ""
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
	if !strings.HasPrefix(resp.Msg, "Please") {
		t.Errorf("Got %s", resp.Msg)
	}
	// Test with bad email address
	resp = response.Info{}
	testUser = Credentials{Name: "Admin", Email: ".com", Pass: "123456"}
	requestBody, err = json.Marshal(testUser)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
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
	requestBody, err = json.Marshal(testUser)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	resp = Register(req, resp)
	got = resp.Msg
	want = "ok"
	if want != got {
		t.Errorf("Got %s want %s", got, want)
	}
}

func TestSessionValid(t *testing.T) {
	files.Root = "../files/"
	// New user
	settings.Set(settings.Values{}, "test")
	got := SessionValid([]string{""}, "test")
	want := true
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Session expired
	state := settings.Values{
		SessionExpiry: time.Now().AddDate(0, 0, -1),
		SessionID:     "aaa",
		Email:         "a@b.co",
	}
	settings.Set(state, "test")
	got = SessionValid([]string{"aaa"}, "test")
	want = false
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Session not expired, wrong id
	state.SessionExpiry = time.Now().AddDate(0, 0, 1)
	settings.Set(state, "test")
	got = SessionValid([]string{"wrongID"}, "test")
	want = false
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Session not expired, correct id
	got = SessionValid([]string{"aaa"}, "test")
	want = true
	if got != want {
		t.Errorf("Got %v want %v", got, want)
	}
	// Clear the settings
	settings.Set(settings.Values{}, "test")
}
