package main

import (
	"alpha/go/gopl/ch4/4_12/elasticsearch"
	"alpha/go/gopl/ch4/4_12/xkcd"
	"flag"
	"fmt"
	"log"
)

var method = flag.String("method", "", "the action, which shall be search/record")
var keyword = flag.String("keyword", "", "the key word to search in the database")
var numBegin = flag.Int64("begin", 600, "the begin(inclusive) number to record")
var numEnd = flag.Int64("end", 600, "the end(exclusive) number to record")

func main() {
	flag.Parse()

	switch *method {
	case "search":
		search()
	case "record":
		record()
	default:
		flag.PrintDefaults()
	}
}

func record() {
	for i := *numBegin; i < *numEnd; i++ {
		item, err := xkcd.Access(int32(i))
		if err != nil {
			log.Fatal(err)
		}
		json, _ := item.JSON()
		fmt.Println(json)
		index, err := elasticsearch.Index(json)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(index)
	}
}

func search() {
	result, err := elasticsearch.Search(*keyword)
	if err != nil {
		log.Fatal(err)
	}
	for i, item := range result.Hits.Hits {
		fmt.Printf(`
No.%-4d score = %v
#%-4d %s %s
%s
`, i+1, item.Score, item.Source.Num, item.Source.Title, item.Source.Img, item.Source.Transcript)
	}
}
