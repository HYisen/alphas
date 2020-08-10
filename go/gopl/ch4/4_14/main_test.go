package main

import "testing"

func TestFindNextFromLink(t *testing.T) {
	tests := []struct {
		mean string
		link string
		next string
	}{
		{"head one", `<https://api.github.com/repositories/136815833/issues?page=2>; rel="next", <https://api.github.com/repositories/136815833/issues?page=3>; rel="last"`, "https://api.github.com/repositories/136815833/issues?page=2"},
		{"medium one", `<https://api.github.com/repositories/20929025/issues?page=1>; rel="prev", <https://api.github.com/repositories/20929025/issues?page=3>; rel="next", <https://api.github.com/repositories/20929025/issues?page=153>; rel="last", <https://api.github.com/repositories/20929025/issues?page=1>; rel="first"`, "https://api.github.com/repositories/20929025/issues?page=3"},
		{"tail one", `<https://api.github.com/repositories/136815833/issues?page=1>; rel="prev", <https://api.github.com/repositories/136815833/issues?page=3>; rel="next", <https://api.github.com/repositories/136815833/issues?page=3>; rel="last", <https://api.github.com/repositories/136815833/issues?page=1>; rel="first"`, ""},
	}

	for _, tc := range tests {
		if actual := findNextFromLink(tc.link); actual != tc.next {
			t.Fatalf("failed %s, expect %s , actual %s", tc.mean, tc.next, actual)
		}
	}
}
