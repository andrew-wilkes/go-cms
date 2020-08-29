package router

import "fmt"

// Request type
type Request struct {
	Domain string
	Route  string
	Params map[string]string
}

// Process a request
func Process(r Request) {
	fmt.Println(r.Route)
}
