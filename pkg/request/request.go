package request

import "encoding/json"

// Info type
type Info struct {
	Domain    string
	Route     string
	SubRoutes []string
	Method    string
	Scheme    string
	GetArgs   map[string]string
	PostData  map[string]json.RawMessage
}
