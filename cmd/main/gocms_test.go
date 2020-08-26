package main

import (
	"io/ioutil"
	"testing"
)

func TestProcessInput(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("test-data-files/args-%2Ftest%3Fa%3D123%26b%3D456.json")
	jsonDataStr := string(jsonData)
	_, _ = ProcessInput(jsonDataStr) // headers, content :=
}
