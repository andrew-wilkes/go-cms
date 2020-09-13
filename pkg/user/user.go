package user

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gocms/pkg/request"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"time"

	"github.com/rs/xid"
)

// Credentials type
type Credentials struct {
	Name  string
	Email string
	Pass  string
	ID    string
}

// Status type
type Status struct {
	UserName string
	Email    string
	LastTime time.Time
}

var cred Credentials

// LogOn to the system
func LogOn(req request.Info, resp response.Info) response.Info {
	ok := true
	err := json.Unmarshal(req.PostData["user"], &cred)
	if err != nil {
		ok = false
		resp.Msg = "Error decoding data!"

	}
	info := settings.Get(req.Domain)
	if ok && info.Email != cred.Email {
		ok = false
		resp.Msg = "Unknown user!"
	}
	if ok && info.Password != hash(cred.Pass) {
		ok = false
		resp.Msg = "Wrong password!"
	}
	if ok {
		// Start a new session
		info.SessionID = xid.New().String()
		info.SessionExpiry = time.Now().Add(time.Hour * 8)
		// Set up the response data
		resp.ID = info.SessionID
		resp.Msg = "ok"
		status := Status{UserName: info.UserName, Email: info.Email, LastTime: info.LastTime}
		data, _ := json.Marshal(status)
		resp.Data = string(data)
		// Save the new info to settings
		info.LastTime = time.Now()
		info.LastIP = req.IPAddr
		settings.Set(info, req.Domain)
	}
	return resp
}

// LogOff from the App
func LogOff(req request.Info, resp response.Info) response.Info {
	info := settings.Get(req.Domain)
	info.SessionID = ""
	info.SessionExpiry = time.Now()
	settings.Set(info, req.Domain)
	resp.Msg = "ok"
	return resp
}

// Register with the App or update user details
func Register(req request.Info, resp response.Info) response.Info {
	info := settings.Get(req.Domain)
	ok := true
	err := json.Unmarshal(req.PostData["user"], &cred)
	if err != nil {
		ok = false
		resp.Msg = "Error decoding data!"
	}
	// We are either setting the details for the first time (cred.ID == "")
	// or updating the user details (cred.ID == current session ID)
	if ok && len(cred.Pass) < 4 {
		ok = false
		resp.Msg = "Please enter a password longer than 3 characters."
	}
	if ok {
		info.UserName = cred.Name
		info.Email = cred.Email
		info.Password = hash(cred.Pass)
		info.SessionID = xid.New().String()
		info.SessionExpiry = time.Now().Add(time.Hour * 8)
		info.LastTime = time.Now()
		info.LastIP = req.IPAddr
		settings.Set(info, req.Domain)
		resp.ID = info.SessionID
		resp.Msg = "ok"
	}
	return resp
}

// Return hashed string
func hash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
