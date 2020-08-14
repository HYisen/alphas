package main

import (
	"alphas/go/gopl/github"
	"alphas/go/gopl/utility"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	milestones, users := parse(issues)

	println("start serving")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer,
			"count\nissues=%d\nmilestones=%d\nusers=%d\n", len(issues), len(milestones), len(users))
	})
	http.HandleFunc("/issue", func(writer http.ResponseWriter, request *http.Request) {
		err = issueTemplate.Execute(writer, struct {
			Count int
			Items []*github.Issue
		}{len(issues), issues})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			writer.WriteHeader(500)
		}
	})
	http.HandleFunc("/milestone", func(writer http.ResponseWriter, request *http.Request) {
		err = milestoneTemplate.Execute(writer, struct {
			Count int
			Items []*github.Milestone
		}{len(issues), milestones})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			writer.WriteHeader(500)
		}
	})
	http.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		err = userTemplate.Execute(writer, struct {
			Count int
			Items []*github.User
		}{len(issues), users})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			writer.WriteHeader(500)
		}
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
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

func parse(issues []*github.Issue) ([]*github.Milestone, []*github.User) {
	var milestones []*github.Milestone
	var users []*github.User
	for _, issue := range issues {
		if issue.User != nil {
			users = append(users, issue.User)
		}
		if issue.Milestone != nil {
			milestones = append(milestones, issue.Milestone)
		}
	}
	return milestones, users
}

var issueTemplate = template.Must(template.New("issue").Parse(`
<h1>{{.Count}} issues</h1>
<table>
    <tr style='text-align: left'>
        <th>#</th>
        <th>State</th>
        <th>User</th>
        <th>Title</th>
    </tr>
    {{range .Items}}
        <tr>
            <td><a href="{{.HTMLURL}}">{{.Number}}</a></td>
            <td>{{.State}}</td>
            <td><a href="{{.User.HTMLURL}}">{{.User.Login}}</a></td>
            <td><a href="{{.HTMLURL}}">{{.Title}}</a></td>
        </tr>
    {{end}}
</table>`))

var milestoneTemplate = template.Must(template.New("milestone").Parse(`
<h1>{{.Count}} milestones</h1>
<table>
    <tr style='text-align: left'>
        <th>#</th>
        <th>State</th>
        <th>Title</th>
        <th>Creator</th>
    </tr>
    {{range .Items}}
        <tr>
            <td><a href="{{.HTMLURL}}">{{.Id}}</a></td>
            <td>{{.State}}</td>
            <td>{{.Title}}</a></td>
            <td><a href="{{.Creator.HTMLURL}}">{{.Creator.Login}}</a></td>
        </tr>
    {{end}}
</table>`))

var userTemplate = template.Must(template.New("user").Parse(`
<h1>{{.Count}} users</h1>
<table>
    <tr style='text-align: left'>
        <th>#</th>
        <th>Name</th>
    </tr>
    {{range .Items}}
        <tr>
            <td><a href="{{.HTMLURL}}">{{.Id}}</a></td>
            <td>{{.Login}}</td>
        </tr>
    {{end}}
</table>`))
