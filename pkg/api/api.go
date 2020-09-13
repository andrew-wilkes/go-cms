package api

import (
	"encoding/json"
	"gocms/pkg/request"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"gocms/pkg/user"
	"time"
)

var data dataPacket
var state settings.Values
var authorized bool

// Process takes action
func Process(req request.Info, domain string) ([]string, string) {
	// Decode the post data for user credentials
	//err := json.Unmarshal(r.PostData["user"], &(data.user))
	if len(r.SubRoutes) > 1 {
		state = settings.Get(domain)
		authorized = sessionValid(domain, req.GetArgs["id"])
		class := r.SubRoutes[0]
		action := r.SubRoutes[1]
		var resp response.Info
		switch class {
		case "user":
			resp = userActions(action, req, resp)
		case "page":
			resp = pageActions(action, req, resp)
		case "pages":
			resp = pagesActions(action, req, resp)
		}
	}
	b, _ := json.Marshal(resp)
	return []string{""}, string(b)
}

func userActions(action string, req request.Info, resp response.Info) response.Info {
	if action == "logon" {
		resp = user.LogOn(req, resp)
	}
	if authorized {
		switch action {
		case "logoff":
			resp = user.LogOff(req, resp)
		case "register":
			resp = user.Register(req, resp)
		}
	}
	return resp
}

func pageActions(action string, req request.Info, resp response.Info) response.Info {
	return resp
}

func pagesActions(action string, req request.Info, resp response.Info) response.Info {
	return resp
}

func sessionValid(domain string, id string) {
	if state.SessionExpiry.Before(time.Now()) {
		authorized = false
	} else {
		authorized = state.SessionID == id
	}
	if authorized {
		response.id = id
	} else {
		response.id = ""
	}
}
