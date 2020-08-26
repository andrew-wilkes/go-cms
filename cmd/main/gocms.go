package main

import (
	"encoding/json"
	"fmt"
	"os"
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

// ProcessInput creates headers and content according to the server input info.
func ProcessInput(jsonData string) ([]string, string) {
	data := []byte(jsonData)
	var serverVars map[string]json.RawMessage // The data values are strings and numbers
	err := json.Unmarshal(data, &serverVars)
	var uriStr string
	var methodStr string
	var rawDataStr string
	if err == nil {
		uri := serverVars["REQUEST_URI"]
		err = json.Unmarshal(uri, &uriStr)
		method := serverVars["REQUEST_METHOD"]
		err = json.Unmarshal(method, &methodStr)
		rawData := serverVars["RAW_DATA"]
		err = json.Unmarshal(rawData, &rawDataStr)
	}

	html := fmt.Sprintf("Err: %v URI: %v Method: %v Data: %v\n", err, uriStr, methodStr, rawDataStr)
	return headers(), content(html)
}

func headers() []string {
	return []string{"apple", "peach", "pear"}
}

func content(html string) string {
	return fmt.Sprintf("%v", html)
}
