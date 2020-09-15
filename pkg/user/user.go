package user

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gocms/pkg/response"
	"gocms/pkg/settings"
	"net/http"
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
func LogOn(req *http.Request, resp response.Info) response.Info {
	ok := true
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&cred)
	if err != nil {
		ok = false
		resp.Msg = "Error decoding data!"

	}
	info := settings.Get(req.Host)
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
		info.LastIP = req.RemoteAddr
		settings.Set(info, req.Host)
	}
	return resp
}

// LogOff from the App
func LogOff(req *http.Request, resp response.Info) response.Info {
	info := settings.Get(req.Host)
	info.SessionID = ""
	info.SessionExpiry = time.Now()
	settings.Set(info, req.Host)
	resp.Msg = "ok"
	return resp
}

// Register with the App or update user details
func Register(req *http.Request, resp response.Info) response.Info {
	info := settings.Get(req.Host)
	ok := true
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&cred)
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
		info.LastIP = req.RemoteAddr
		settings.Set(info, req.Host)
		resp.ID = info.SessionID
		resp.Msg = "ok"
	}
	return resp
}

// Return hashed string
func hash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
