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

func main() {
	args := os.Args[1:]
	getArg := ""
	if len(args) > 0 {
		getArg = args[0]
	}
	postArg := ""
	if len(args) > 1 {
		postArg = args[1]
	}
	resA := &response{
		Headers: []string{"apple", "peach", "pear"},
		Content: fmt.Sprintf("The content. Get: %s Post: %s", getArg, postArg),
	}
	resB, _ := json.Marshal(resA)
	fmt.Print(string(resB))
}
