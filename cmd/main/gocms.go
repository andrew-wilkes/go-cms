package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// URI contains the route and GET parameters
var URI string

// Method is the request method (GET or POST)
var Method string

// RawData contains any json POST data
var RawData string

// Domain is used for the file path to data and pages
var Domain string

type response struct {
	Headers []string
	Content string
}

// Request type
type Request struct {
	route  string
	params map[string]string
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

// ProcessInput creates headers and content according to the server input info.
func ProcessInput(jsonData string) ([]string, string) {
	data := []byte(jsonData)
	var serverVars map[string]json.RawMessage // The data values are strings and numbers
	err := json.Unmarshal(data, &serverVars)
	if err == nil {
		uri := serverVars["REQUEST_URI"]
		err = json.Unmarshal(uri, &URI)
		method := serverVars["REQUEST_METHOD"]
		err = json.Unmarshal(method, &Method)
		rawData := serverVars["RAW_DATA"]
		err = json.Unmarshal(rawData, &RawData)
		domain := serverVars["SERVER_NAME"]
		err = json.Unmarshal(domain, &Domain)
	}

	html := fmt.Sprintf("Error: %v URI: %v Method: %v Data: %v Domain: %v\n", err, URI, Method, RawData, Domain)
	return headers(), content(html)
}

func headers() []string {
	return []string{"apple", "peach", "pear"}
}

func content(html string) string {
	return fmt.Sprintf("%v", html)
}

// ParseURI splits the uri into component parts
func ParseURI(uri string) Request {
	var r = Request{params: make(map[string]string)}
	p := strings.Split(uri, "?")
	r.route = p[0]
	if len(p) > 1 {
		params := strings.Split(p[1], "&")
		for _, param := range params {
			kv := strings.Split(param, "=")
			r.params[kv[0]] = kv[1]
		}
	}
	return r
}
