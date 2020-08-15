package utility

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
)

func Fetch(method string, url string, jsonStringPayload io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, jsonStringPayload)
	if err != nil {
		return nil, err
	}
	if jsonStringPayload != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

func GetHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer CloseAndLogError(resp.Body)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
