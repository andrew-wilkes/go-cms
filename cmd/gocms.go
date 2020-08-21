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
	info := os.Args[1:]

	headers, content := ProcessInput(info)

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
func ProcessInput(info []string) ([]string, string) {
	// The info will be comprised of a map of key:values
	return headers(), content(info)
}

func headers() []string {
	return []string{"apple", "peach", "pear"}
}

func content(info []string) string {
	return fmt.Sprintf("The server info: %v", info)
}
