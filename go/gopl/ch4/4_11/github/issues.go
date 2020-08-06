package github

import "time"

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
