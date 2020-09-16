package main

import (
	"fmt"
	"gocms/pkg/router"
	"net/http"
)

// This is the entry point function for the application
func main() {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8090", nil)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	statusCode, headers, content := router.Process(req)
	for k, v := range headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, content)
}
