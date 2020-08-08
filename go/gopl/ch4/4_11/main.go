package main

import (
	"alpha/go/gopl/ch4/4_11/github"
	"alpha/go/gopl/ch4/4_11/utility"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var path = flag.String("path", "", "whose issues is going to be managed")
var method = flag.String("method", "", "action on issues : create/read/update/delete")
var number = flag.Int("number", -1, "detail issue number to be handled")
var token = flag.String("token", "", "GitHub Access Token for privileges")

var url string

func main() {
	flag.Parse()
	url = fmt.Sprintf("https://api.github.com/repos/%s/issues", *path)

	switch *method {
	case "read":
		issues, err := Retrieve()
		if err != nil {
			log.Panic(err)
		}

		if *number == -1 {
			for _, issue := range issues {
				//fmt.Println(issue)
				fmt.Printf("#%-6d %16s %16s %s %s\n", issue.Number, issue.Title,
					issue.User.Login, issue.CreatedAt.Format("060102"), issue.UpdatedAt.Format("060102"))
			}
		} else {
			if issue, ok := findIssueById(issues, int32(*number)); ok {
				fmt.Printf("#%d %s \n  by %s\n  created at %v\n  updated at %v\n%s\n",
					issue.Number, issue.Title, issue.User.Login, issue.CreatedAt, issue.UpdatedAt, issue.Body)
			} else {
				log.Panicf("can not found issue #%d", *number)
			}
		}
	case "create":
		issue := GenNeoIssue()
		fmt.Println(issue)
		Create(issue)
	default:
		flag.PrintDefaults()
	}
}

func findIssueById(array []github.Issue, target int32) (github.Issue, bool) {
	for _, issue := range array {
		if issue.Number == target {
			return issue, true
		}
	}
	return github.Issue{}, false
}

func Retrieve() ([]github.Issue, error) {
	resp, err := fetch("GET", nil)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("bad resp status %v", resp.Status)
	}

	var issues []github.Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("can not decode : %v", err)
	}

	_ = resp.Body.Close()
	return issues, nil
}

func GenNeoIssue() github.NeoIssue {
	const msg = `
# Please enter the issue title and body, which shall be separated by
# the first empty line, without which means an issue has no body but
# only title. Lines starting with '#' will be ignored, and an empty
# message aborts the commit.
#`

	input := utility.GetInputFromTextEditor("echo", msg)
	// Null or blank body are identical to GitHub web issue detail page.
	// Therefore the distinguish there are just for fun
	title, body, isBodyExist := utility.ExtractTitleAndBody(input)

	neo := github.NeoIssue{
		Title: title,
	}
	fmt.Println(isBodyExist)
	if isBodyExist {
		neo.Body = &body
	}

	return neo
}

func fetch(method string, payload io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if len(*token) != 0 {
		req.Header.Add("Authorization", " token "+*token)
	}
	return http.DefaultClient.Do(req)
}

func Create(neoIssue github.NeoIssue) {
	var sb strings.Builder
	if err := json.NewEncoder(&sb).Encode(neoIssue); err != nil {
		log.Fatal(err)
	}

	resp, err := fetch("POST", strings.NewReader(sb.String()))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		log.Fatalf("bad return http status %v", resp.Status)
	}
	fmt.Println(resp.Body)
	_ = resp.Body.Close()
}
