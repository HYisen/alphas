package xkcd

import (
	"alphas/go/gopl/utility"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// I don't want any redundant data for future usage. As YAGNI says.
type Item struct {
	//Year       string `json:"year"`
	//Month      string `json:"month"`
	//Day        string `json:"day"`
	Num int32 `json:"num"`
	//Link       string `json:"link"`
	//News       string `json:"news"`
	Title string `json:"title"`
	//SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	//Alt        string `json:"alt"`
	Img string `json:"img"`
}

func Access(num int32) (*Item, error) {
	resp, err := http.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", num))
	if err != nil {
		return nil, err
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response status %v", resp.Status)
	}

	var item Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, err
}

func (i Item) JSON() (string, error) {
	sb := strings.Builder{}
	err := json.NewEncoder(&sb).Encode(i)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
