package main

import (
	"encoding/json"
	"fmt"
	"gocms/pkg/router"
	"os"
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
	// Get the passed array of information ignoring the first parameter which is a file name
	// https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	// https://www.geeksforgeeks.org/how-to-parse-json-in-golang/

	headers, content := ProcessInput(os.Args[1:][0])

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

// uri contains the route and GET parameters
var uri string

// method is the request method (GET or POST)
var method string

// rawData contains any json POST data
var rawData string

// domain is used for the file path to data and pages
var domain string

// ProcessInput creates headers and content according to the server input info.
func ProcessInput(jsonData string) ([]string, string) {
	serverVars, err := DecodeRawData(jsonData)
	if err == nil {
		_uri := serverVars["REQUEST_URI"]
		err = json.Unmarshal(_uri, &uri)
		_method := serverVars["REQUEST_METHOD"]
		err = json.Unmarshal(_method, &method)
		_rawData := serverVars["RAW_DATA"]
		err = json.Unmarshal(_rawData, &rawData)
		_domain := serverVars["SERVER_NAME"]
		err = json.Unmarshal(_domain, &domain)
		r := ParseURI(uri)
		r.Domain = domain
		router.Process(r)
	}

	html := fmt.Sprintf("Error: %v uri: %v Method: %v Data: %v Domain: %v\n", err, uri, method, rawData, domain)
	return headers(), content(html)
}

func headers() []string {
	return []string{"apple", "peach", "pear"}
}

func content(html string) string {
	return fmt.Sprintf("%v", html)
}

// ParseURI splits the uri into component parts
func ParseURI(uri string) router.Request {
	var r = router.Request{Params: make(map[string]string)}
	p := strings.Split(uri, "?")
	r.Route = p[0]
	if len(p) > 1 {
		params := strings.Split(p[1], "&")
		for _, param := range params {
			kv := strings.Split(param, "=")
			r.Params[kv[0]] = kv[1]
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
