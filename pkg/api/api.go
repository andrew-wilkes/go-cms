package api

import (
	"encoding/json"
	"gocms/pkg/page"
	"gocms/pkg/request"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"gocms/pkg/user"
	"time"
)

var state settings.Values
var authorized bool

// Process takes action
func Process(req request.Info) ([]string, string) {
	var resp response.Info
	if len(req.SubRoutes) > 1 {
		state = settings.Get(req.Domain)
		sessionValid(req.Domain, req.GetArgs["id"])
		class := req.SubRoutes[0]
		action := req.SubRoutes[1]
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
	// Add actions related to saving edited page content from the in-page editor (Content Tools)
	if authorized {
		switch action {
		case "save":
			var info page.EditInfo
			err := json.Unmarshal(req.PostData["info"], &info)
			if err != nil {
				resp.Msg = "Error decoding data!"
			} else {
				page.SaveContent(req.Domain, info.ID, info.Content)
				resp.Msg = "ok"
			}
		}
	}
	return resp
}

func pagesActions(action string, req request.Info, resp response.Info) response.Info {
	// Add actions related to save and load of pages data from the Dashboard
	if authorized {
		switch action {
		case "save":
			page.SaveRawData(req.Domain, []byte(req.PostData["pages"]))
		case "load":
			resp.Data = string(page.LoadRawData(req.Domain))
		}
		resp.Msg = "ok"
	}
	return resp
}

func sessionValid(domain string, id string) {
	if state.SessionExpiry.Before(time.Now()) {
		authorized = false
	} else {
		authorized = state.SessionID == id
	}
}
