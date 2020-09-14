package main

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
	}{
		{name: "Basic", s: "aba", sep: "b"},
		{name: "LongerSep", s: "mea in the research lab", sep: "ea"},
		{name: "UTF", s: "もっとaaもっとお喋りしていたいaaんです", sep: "い"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := split(test.s, test.sep)
			expect := strings.Split(test.s, test.sep)
			if !equal(actual, expect) {
				t.Fatalf("faild on split %s with %s, expect %s, but actual %s", test.s, test.sep, expect, actual)
			}
		})
	}
}

func equal(one, two []string) bool {
	if len(one) != len(two) {
		return false
	}
	for i := range one {
		if one[i] != two[i] {
			return false
		}
	}
	return true
}
