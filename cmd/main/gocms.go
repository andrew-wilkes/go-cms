package main

import (
	"fmt"
	"gocms/pkg/router"
	"net/http"
)

// This is the entry point function for the application
func main() {
	const port = 8090
	fmt.Printf("Listening on port: %d\n", port)
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	statusCode, headers, content := router.Process(req)
	for k, v := range headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, content)
}
