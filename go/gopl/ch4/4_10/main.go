package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// issues

type User struct {
	Login   string
	Id      int
	HTMLURL string `json:"html_url"`
}

type Issue struct {
	Id        int
	HTMLURL   string `json:"html_url"`
	Number    int
	State     string
	Title     string
	Body      string
	User      *User
	CreatedAt time.Time `json:"created_at"`
}

type IssueSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

const IssuesURL = "https://api.github.com/search/issues"

//const IssuesURL="https://api.github.com/dmhacker/arch-linux-surface/issues"

func SearchIssues(term []string) (*IssueSearchResult, error) {
	q := url.QueryEscape(strings.Join(term, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("bad resp status %v", resp.Status)
	}

	var result IssueSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	_ = resp.Body.Close()
	return &result, nil
}

func main() {
	args := [...]string{"repo:golang/go", "is:open", "json", "decoder"}
	issues, err := SearchIssues(args[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count=%d\n", issues.TotalCount)
	items := issues.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	hint := [...]string{"less than a month", "less than a year old", "more than a year old"}
	now := time.Now()
	checkIsExceededMethods := [...]func(data time.Time) bool{
		func(data time.Time) bool {
			return data.Add(time.Hour * 24 * 30).Before(now)
		},
		func(data time.Time) bool {
			return data.Add(time.Hour * 24 * 365).Before(now)
		},
		func(data time.Time) bool {
			return false
		},
	}

	stage := 0
	fmt.Printf("\n%s\n", hint[stage])
	for _, item := range items {
		if checkIsExceededMethods[stage](item.CreatedAt) {
			stage++
			fmt.Printf("\n%s\n", hint[stage])
		}
		fmt.Printf("#%-5d %v %9.9s %.55s\n", item.Number, item.CreatedAt, item.User.Login, item.Title)
	}
}
