package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_extractArg(t *testing.T) {
	main.args := []string{"a", "b"}
	ans := main.extractArg(0)
    if ans != "a" {
		t.Errorf("main.extractArg(0) = %s; want 'a'", ans)
	}
}
