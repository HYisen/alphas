package main

import (
	"alphas/go/gopl/utility"
	"fmt"
	"golang.org/x/net/html"
	"log"
)

func main() {
	doc, err := utility.GetHTML("https://golang.org")
	if err != nil {
		log.Fatal(err)
	}
	nodes := ElementsByTagName(doc, "title", "h1", "h2")
	for _, node := range nodes {
		fmt.Printf("%s : %s\n", node.Data, node.FirstChild.Data)
	}
}

// doc shall not be nil.
func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	// build set
	set := make(map[string]bool)
	for _, item := range name {
		set[item] = true
	}

	var result []*html.Node

	candidates := []*html.Node{doc}
	for len(candidates) != 0 {
		node := candidates[0]
		candidates = candidates[1:] // pop first
		if node.Type == html.ElementNode {
			if _, ok := set[node.Data]; ok {
				result = append(result, node)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			candidates = append(candidates, child) // push last
		}
	}

	return result
}
