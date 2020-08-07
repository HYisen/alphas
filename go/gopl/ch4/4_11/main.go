package main

import (
	"alpha/go/gopl/ch4/4_11/github"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var path = flag.String("path", "", "whose issues is going to be managed")
var method = flag.String("method", "", "action on issues : create/read/update/delete")
var number = flag.Int("number", -1, "detail issue number to be handled")

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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)

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
