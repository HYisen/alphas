package main

import (
	"alphas/go/gopl/utility"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Result struct {
	Title  string `json:"Title"`
	Poster string `json:"Poster"`

	Response string `json:"Response"`
	Error    string `json:"Error"`
}

var title = flag.String("title", "", "title of the movie")
var apikey = flag.String("apikey", "", "apikey of OMDb API")

func main() {
	flag.Parse()

	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&apikey=%s", *title, *apikey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer utility.CloseAndLogError(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatalf("bad status %v", resp.Status)
	}

	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Response == "False" {
		log.Fatalf("failed to search : %v", result.Error)
	}

	more, err := http.Get(result.Poster)
	if err != nil {
		log.Fatal(err)
	}
	defer utility.CloseAndLogError(more.Body)

	if more.StatusCode != 200 {
		log.Fatalf("bad status %v", more.Status)
	}

	ext := result.Poster[strings.LastIndex(result.Poster, "."):]
	file, err := os.Create(result.Title + ext)
	if err != nil {
		log.Fatal(err)
	}
	defer utility.CloseAndLogError(file)

	_, err = io.Copy(file, more.Body)
	if err != nil {
		log.Fatal(err)
	}
}
