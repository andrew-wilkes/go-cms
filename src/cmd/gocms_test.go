package main

import (
	"testing"
)

func Test_ExtractArg(t *testing.T) {
	var tests = []struct {
		args []string
		a    string
		b    string
	}{
		{[]string{}, "", ""},
		{[]string{"a"}, "a", ""},
		{[]string{"a", "b"}, "a", "b"},
	}
	for _, tt := range tests {
		aa := ExtractArg(tt.args, 0)
		bb := ExtractArg(tt.args, 1)
		if aa != tt.a || bb != tt.b {
			t.Errorf("ExtractArg(args, index); want 'arg[0] == %s and arg[1] == %s'; got 'arg[0] == %s and arg[1] == %s'", tt.a, tt.b, aa, bb)
		}
	}
}
