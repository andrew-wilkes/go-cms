package main

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/files"
	"gocms/pkg/request"
	"gocms/pkg/router"
	"os"
	"path/filepath"
	"strings"
)

type response struct {
	Headers []string
	Content string
}

// This is the entry point function for the application
// It is called (from index.php) with command line arguments representing
// the source file name and an array of information passed from the server about the request
// See: https://www.php.net/manual/en/reserved.variables.server.php
func main() {
	files.Root, _ = filepath.Split(os.Args[0])
	headers, content := ProcessInput(os.Args[1])

	// Construct the response containing headers and content
	resA := &response{
		Headers: headers,
		Content: content,
	}
	// Encode the response in text format
	resB, _ := json.Marshal(resA)

	// Output the response
	fmt.Print(string(resB))
}

// ProcessInput creates headers and content according to the server input info.
func ProcessInput(jsonData string) ([]string, string) {
	headers := []string{}
	content := ""
	serverVars, err := DecodeRawData(jsonData)
	if err == nil {
		var uri string // contains the route and GET parameters
		_uri := serverVars["REQUEST_URI"]
		err = json.Unmarshal(_uri, &uri)
		r := ParseURI(uri)
		_domain := serverVars["SERVER_NAME"] // this is used for the file path to data and pages
		err = json.Unmarshal(_domain, &r.Domain)
		_method := serverVars["REQUEST_METHOD"] // the request method (GET or POST)
		err = json.Unmarshal(_method, &r.Method)
		_scheme := serverVars["REQUEST_SCHEME"] // http or https
		err = json.Unmarshal(_scheme, &r.Scheme)
		_ipAddr := serverVars["REMOTE_ADDR"]
		err = json.Unmarshal(_ipAddr, &r.IPAddr)
		var rawData string // contains any json POST data
		_rawData := serverVars["RAW_DATA"]
		err = json.Unmarshal(_rawData, &rawData)
		// Decode the json string
		r.PostData, err = DecodeRawData(rawData)
		headers, content = router.Process(r)
	} else {
		content = "Error decoding input data!"
	}
	return headers, content
}

// ParseURI splits the uri into component parts
func ParseURI(uri string) request.Info {
	var r = request.Info{GetArgs: make(map[string]string)}
	p := strings.Split(uri, "?")
	r.Route = p[0]
	if len(p) > 1 {
		args := strings.Split(p[1], "&")
		for _, arg := range args {
			kv := strings.Split(arg, "=")
			r.GetArgs[kv[0]] = kv[1]
		}
	}
	return r
}

// DecodeRawData decodes a json text string to a map of values
func DecodeRawData(rawData string) (map[string]json.RawMessage, error) {
	data := []byte(rawData)
	var keyValues map[string]json.RawMessage // The data values are strings and numbers
	err := json.Unmarshal(data, &keyValues)
	return keyValues, err
}
