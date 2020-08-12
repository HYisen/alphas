package main

import (
	"alpha/go/gopl/utility"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	url := "https://hyisen.net/"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer utility.CloseAndLogError(resp.Body)

	var sb strings.Builder
	_, _ = io.Copy(&sb, resp.Body)
	fmt.Println(sb.String())

	fmt.Println("\nHere We Go!")

	doc, err := html.Parse(strings.NewReader(sb.String()))
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		log.Fatal(err)
	}

	elem := ElementByID(doc, "root")
	fmt.Println(elem == nil)
	if elem != nil {
		fmt.Println(elem.Data)
		fmt.Println(elem.Attr)
	}
}

func ElementByID(doc *html.Node, id string) *html.Node {
	keyword = id
	forEachNode(doc, startElement, nil)
	return target
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) {
	if pre != nil {
		if pre(n) {
			return
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		if post(n) {
			return
		}
	}
}

var target *html.Node
var keyword string

func startElement(n *html.Node) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" {
				if attr.Val == keyword {
					target = n
					return false
				}
				return true
			}
		}
	}
	return true
}
