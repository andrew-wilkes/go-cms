package api

import (
	"encoding/json"
	"gocms/pkg/request"
	"gocms/pkg/user"
	"testing"
)

type postData struct {
	user user.Credentials
}

func TestProcess(t *testing.T) {
	testUser := user.Credentials{"Admin", "a@b.com", "pwd", "xyz"}
	data, _ := json.Marshal(testUser)
	r := request.Info{
		PostData: map[string]json.RawMessage{"user": data},
	}
	_, got := Process(r)
	want := ""
	if got != want {
		t.Errorf("Want %v got %v", want, got)
	}
}
