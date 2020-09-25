package settings

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"io/ioutil"
	"path/filepath"
	"time"
)

// Values type
type Values struct {
	loaded        bool
	UserName      string
	Email         string
	Password      string
	LastIP        string
	LastTime      time.Time
	SessionID     string
	SessionExpiry time.Time
}

var values Values

// Get values
func Get(domain string) Values {
	if values.loaded == false {
		load(domain)
		values.loaded = true
	}
	return values
}

// Set values
func Set(v Values, domain string) {
	values = v
	save(domain)
}

func load(domain string) {
	data, err := ioutil.ReadFile(fileName(domain))
	if err == nil {
		err = json.Unmarshal(data, &values)
		if err != nil {
			panic(err)
		}
	} else {
		// Assume that the file doesn't exist so try to create it
		// Set default values here
		save(domain)
	}
}

func save(domain string) {
	b, _ := json.Marshal(values)
	err := ioutil.WriteFile(fileName(domain), b, 0660)
	if err != nil {
		panic(err)
	}
}

func fileName(domain string) string {
	return filepath.Join(fmt.Sprintf("%s%s", files.Root, domain), "data", "settings.json")
}
