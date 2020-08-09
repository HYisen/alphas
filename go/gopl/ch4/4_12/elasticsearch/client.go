package elasticsearch

import (
	"alpha/go/gopl/utility"
	json2 "encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	addr      = "localhost:9200"
	indexName = "xkcd"
)

func Index(json string) (*RunResult, error) {
	url := fmt.Sprintf("http://%s/%s/_doc/", addr, indexName)
	resp, err := http.Post(url, "application/json", strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("bad response status %v", resp.Status)
	}

	fmt.Println(resp.Body)
	var rr RunResult
	err = json2.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}
	return &rr, nil
}

type ShardInfo struct {
	Total      int32 `json:"total"`
	Successful int32 `json:"successful"`
	Failed     int32 `json:"failed"`
}

type RunResult struct {
	Index       string    `json:"_index"`
	Type        string    `json:"_type"`
	Id          string    `json:"_id"`
	Version     int32     `json:"_version"`
	Result      string    `json:"result"`
	Shards      ShardInfo `json:"_shards"`
	SeqNO       int32     `json:"_seq_no"`
	PrimaryTerm int32     `json:"_primary_term"`
}
