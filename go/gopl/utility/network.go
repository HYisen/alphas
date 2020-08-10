package utility

import (
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
