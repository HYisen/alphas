package main

import (
	"testing"
)

func TestExtractWebsite(t *testing.T) {
	*vault = "/tmp/web"
	website = "https://web.net/"
	*first = "https://web.net"

	tests := []struct {
		name        string
		input       string
		expectError bool
		expect      string
	}{
		{name: "UnchangedSlash", input: "https://one.net/", expectError: false, expect: "https://one.net/"},
		{name: "ExtractIndexPage", input: "https://one.net/index", expectError: false, expect: "https://one.net/"},
		{name: "ExtractIndexPageWithExt", input: "http://a.net/index.htm", expectError: false, expect: "http://a.net/"},
		{name: "MalformedIndexURL", input: "https://one.net", expectError: true, expect: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := extractWebsite(test.input)
			if test.expectError {
				if err == nil {
					t.Fatalf("expect error but none")
				}
			}
			if err != nil && !test.expectError {
				t.Fatalf("face unexpected error %v", err)
			}
			if actual != test.expect {
				t.Fatalf("expect %s but actual %s\n", test.expect, actual)
			}
		})
	}
}

func TestGenNeoHref(t *testing.T) {
	*vault = "/tmp/web"
	website = "https://web.net/"
	*first = "https://web.net/" // malformed IndexURL shall have been excluded by extractWebsite

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "AddExt", input: "file", expect: "./file.html"},
		{name: "HandleProtocol", input: "https://web.net/file.html", expect: "./file.html"},
		{name: "UnchangedOnDot", input: "./file.htm", expect: ""},
		{name: "UnchangedOnDotDot", input: "../../file.htm", expect: ""},
		{name: "ModifyNoneIndexPage", input: "https://web.net", expect: "./index.html"},
		{name: "ModifySlashIndexPage", input: "https://web.net/", expect: "./index.html"},
		{name: "WithDepthWithoutExt", input: "https://web.net/a/b/c", expect: "./a/b/c.html"},
		{name: "WithDepthWithExt", input: "https://web.net/a/b/c.htm", expect: "./a/b/c.htm"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := genNeoHref(test.input)
			if actual != test.expect {
				t.Fatalf("expect %s but actual %s\n", test.expect, actual)
			}
		})
	}
}

func TestGenFilepath(t *testing.T) {
	*vault = "/tmp/web"
	website = "https://web.net/"
	//*first = "https://web.net/" // malformed IndexURL shall have been excluded by extractWebsite

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "AddExt", input: "./file", expect: "/tmp/web/file.html"},
		{name: "AddExt", input: "file", expect: "/tmp/web/file.html"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := genFilepath(test.input)
			if actual != test.expect {
				t.Fatalf("expect %s but actual %s\n", test.expect, actual)
			}
		})
	}

}
