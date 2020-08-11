package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	traverse(doc)

	fmt.Println(data)
}

var data = map[string]int{}

func traverse(n *html.Node) {
	if n.Type == html.ElementNode {
		data[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
	}
}
