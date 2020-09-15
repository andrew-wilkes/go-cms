package api

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/page"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"gocms/pkg/user"
	"net/http"
	"strconv"
	"time"
)

var state settings.Values
var authorized bool

// Process takes action
func Process(req *http.Request, subRoutes []string) ([]string, string) {
	var resp response.Info
	if len(subRoutes) > 1 {
		state = settings.Get(req.Domain)
		sessionValid(req.GetArgs["id"]) // If the key doesn't exist, then "" is passed maybe?
		class := subRoutes[0]
		action := subRoutes[1]
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

func userActions(action string, req *http.Request, resp response.Info) response.Info {
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

func pageActions(action string, req *http.Request, resp response.Info) response.Info {
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
		// NOTE: This action is unnecessary with in-content editing
		case "load":
			i, er := strconv.Atoi(req.GetArgs["pid"])
			if er == nil {
				info := page.GetByID(req.Domain, i, true)
				if info.ID == 0 {
					resp.Msg = "Page not found!"
				} else {
					resp.Data = info.Content
					resp.Msg = "ok"
				}
			} else {
				resp.Msg = fmt.Sprint(er)
			}
		}
	}
	return resp
}

func pagesActions(action string, req *http.Request, resp response.Info) response.Info {
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

func sessionValid(id string) {
	if state.SessionExpiry.Before(time.Now()) {
		authorized = false
	} else {
		authorized = state.SessionID == id
	}
}
