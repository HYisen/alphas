package main

import (
	"alpha/go/gopl/github"
	"alpha/go/gopl/utility"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var path = flag.String("path", "Dreamacro/clash", "github path of project. {org}/{project}")
var url string

func main() {
	flag.Parse()
	url = fmt.Sprintf("https://api.github.com/repos/%s/", *path)
	fmt.Println(url)

	issues, err := acquireIssues()
	if err != nil {
		log.Fatal(err)
	}
	for _, pi := range issues {
		//fmt.Println(pi.Title)
		fmt.Println(pi.Number)
	}
}

func findNextFromLink(src string) string {
	limbs := strings.Split(src, ",")
	var next string
	var last string
	for _, limb := range limbs {
		limb = strings.TrimSpace(limb)
		link := limb[1:strings.Index(limb, ">")]
		rel := limb[len(limb)-5 : len(limb)-1]
		switch rel {
		case "next":
			next = link
		case "last":
			last = link
		}
	}
	if next == last {
		return ""
	}
	return next
}

func fetch(url string) (string, []*github.Issue, error) {
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		return "", nil, err
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("bad resp status %v", resp.Status)
	}

	var issues []*github.Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return "", nil, fmt.Errorf("can not decode : %v", err)
	}

	return resp.Header.Get("link"), issues, nil
}

func acquireIssues() ([]*github.Issue, error) {
	var issues []*github.Issue

	next := url + "issues"
	for len(next) != 0 {
		link, one, err := fetch(next)
		if err != nil {
			return nil, err
		}
		issues = append(issues, one...)
		next = findNextFromLink(link)
	}

	return issues, nil
}
