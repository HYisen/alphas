package github

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type User struct {
	Id      int64
	Login   string
	HTMLURL string `json:"html_url"`
}

type Issue struct {
	Id        int64
	Number    int32
	Title     string
	Body      string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User
	Milestone *Milestone
}

type NeoIssue struct {
	Title     string   `json:"title"`
	Body      *string  `json:"body,omitempty"`
	Milestone *int32   `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

type Milestone struct {
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	LabelsURL    string    `json:"labels_url"`
	Id           int64     `json:"id"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      *User     `json:"creator"`
	OpenIssues   int32     `json:"open_issues"`
	ClosedIssues int32     `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}

func (i *NeoIssue) JSONReader() *strings.Reader {
	var sb strings.Builder
	if err := json.NewEncoder(&sb).Encode(i); err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(sb.String())
	return reader
}
