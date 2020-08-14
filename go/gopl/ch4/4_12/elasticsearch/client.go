package elasticsearch

import (
	"alphas/go/gopl/ch4/4_12/xkcd"
	"alphas/go/gopl/utility"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	addr      = "localhost:9200"
	indexName = "xkcd"
)

func Index(itemJSON string) (*RunResult, error) {
	url := fmt.Sprintf("http://%s/%s/_doc/", addr, indexName)
	resp, err := http.Post(url, "application/json", strings.NewReader(itemJSON))
	if err != nil {
		return nil, err
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("bad response status %v", resp.Status)
	}

	var rr RunResult
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}
	return &rr, nil
}

type ShardInfo struct {
	Total      int32 `json:"total"`
	Successful int32 `json:"successful"`
	Skipped    int32 `json:"skipped,omitempty"` // would not return in create action
	Failed     int32 `json:"failed"`
}

type RunResult struct {
	basicInfo
	Version     int32     `json:"_version"`
	Result      string    `json:"result"`
	Shards      ShardInfo `json:"_shards"`
	SeqNO       int32     `json:"_seq_no"`
	PrimaryTerm int32     `json:"_primary_term"`
}

type basicInfo struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	Id    string `json:"_id"`
}

type ResultItem struct {
	basicInfo
	Score  float64   `json:"_score"`
	Source xkcd.Item `json:"_source"`
}

type HitInfo struct {
	Value    int32  `json:"value"`
	Relation string `json:"relation"`
}

type SearchResult struct {
	Took     int32     `json:"took"`
	TimedOut bool      `json:"timed_out"`
	Shards   ShardInfo `json:"_shards"`
	Hits     struct {
		Total *HitInfo      `json:"total"`
		Hits  []*ResultItem `json:"hits"`
	} `json:"hits"`
}

func Search(keyword string) (*SearchResult, error) {
	var sr SearchResult

	url := fmt.Sprintf("http://%s/%s/_search", addr, indexName)
	const template = `{
    "query": {
      "multi_match": {
        "query": "%s",
        "fields": ["title","transcript"],
        "fuzziness": 2
      }
    }
  }`

	resp, err := utility.Fetch("GET", url, strings.NewReader(fmt.Sprintf(template, keyword)))
	if err != nil {
		return nil, err
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad resp status %v", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}
