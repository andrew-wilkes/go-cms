package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gocms/pkg/page"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"gocms/pkg/user"
	"io/ioutil"
	"net/http"
	"strconv"
)

var state settings.Values
var authorized bool

// Process takes action
func Process(req *http.Request, subRoutes []string, headers map[string]string) (int, map[string]string, string) {
	var resp response.Info
	// Respond with a message to "register" if no user email has been saved. Or "logon", "ok", or "" on error
	if settings.Get(req.Host).Email == "" {
		resp.Msg = "register"
	} else {
		resp.Msg = "logon"
	}
	if len(subRoutes) > 1 {
		state = settings.Get(req.Host)
		authorized = user.SessionValid(req.URL.Query()["id"], req.Host)
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
	respBytes, _ := json.Marshal(resp)
	return http.StatusOK, headers, string(respBytes)
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
			decoder := json.NewDecoder(req.Body)
			err := decoder.Decode(&info)
			if err != nil {
				resp.Msg = "Error decoding data!"
			} else {
				page.SaveContent(req.Host, info.ID, info.Content)
				resp.Msg = "ok"
			}
		// NOTE: This action is unnecessary with in-content editing
		case "load":
			pid := req.URL.Query()["pid"]
			var err error
			var i int
			if len(pid) != 1 {
				err = errors.New("Missing pid")
			} else {
				i, err = strconv.Atoi(pid[0])
			}
			if err == nil {
				info := page.GetByID(req.Host, i, true)
				if info.ID == 0 {
					resp.Msg = "Page not found!"
				} else {
					resp.Data = info.Content
					resp.Msg = "ok"
				}
			} else {
				resp.Msg = fmt.Sprint(err)
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
			data, _ := ioutil.ReadAll(req.Body)
			page.SaveRawData(req.Host, data)
		case "load":
			resp.Data = string(page.LoadRawData(req.Host))
		}
		resp.Msg = "ok"
	}
	return resp
}
