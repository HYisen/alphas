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
}

type NeoIssue struct {
	Title     string   `json:"title"`
	Body      *string  `json:"body,omitempty"`
	Milestone *int32   `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

func (i *NeoIssue) JSONReader() *strings.Reader {
	var sb strings.Builder
	if err := json.NewEncoder(&sb).Encode(i); err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(sb.String())
	return reader
}
