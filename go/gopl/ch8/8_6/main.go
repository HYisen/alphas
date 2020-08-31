package main

import (
	"alphas/go/gopl/utility"
	"flag"
	"fmt"
	"gopl.io/ch5/links"
	"log"
)

var depth = flag.Int("depth", 2, "the max range to crawl")

type task struct {
	url   string
	depth int
}

func enrich(source []string, depth int) []task {
	var ret []task
	for _, url := range source {
		ret = append(ret, task{
			url:   url,
			depth: depth,
		})
	}
	return ret
}

func main() {
	flag.Parse()

	worklist := make(chan []task)
	unseenLinks := make(chan task)

	go func() { worklist <- enrich(utility.AcquireNoFlagArgs(), 0) }()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link.url)
				go func(neoDepth int) { worklist <- enrich(foundLinks, neoDepth) }(link.depth + 1)
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link.url] {
				seen[link.url] = true
				if link.depth <= *depth {
					unseenLinks <- link
				}
			}
		}
		fmt.Printf("done %v\n", list)
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
