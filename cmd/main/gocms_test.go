package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestProcessInput(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("test-data-files/args-%2Ftest%3Fa%3D123%26b%3D456.json")
	jsonDataStr := string(jsonData)
	_, _ = ProcessInput(jsonDataStr) // headers, content :=
}

func TestParseURI(t *testing.T) {
	r := ParseURI("test/page?a=X&b=Y&c=Z")
	if r.Route != "test/page" && r.Params["a"] != "X" && r.Params["b"] != "Y" && r.Params["c"] != "Z" {
		t.Errorf("Unexpected result in: %v", r)
	}
}

func TestDecodeRawData(t *testing.T) {
	data, err := DecodeRawData(`{ "a": "A", "n": 2 }`)
	if err != nil {
		t.Errorf("%v", err)
	}
	var s string
	var n int
	err = json.Unmarshal(data["a"], &s)
	if err != nil {
		t.Errorf("%v", err)
	}
	err = json.Unmarshal(data["n"], &n)
	if err != nil {
		t.Errorf("%v", err)
	}
	if s != "A" || n != 2 {
		t.Errorf("Bad data: %v %v", s, n)
	}
}
