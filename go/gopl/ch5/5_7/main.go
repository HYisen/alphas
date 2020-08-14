package main

import (
	"alphas/go/gopl/utility"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := "https://hyisen.net"

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

	writer = os.Stdout
	forEachNode(doc, startElement, endElement)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

var depth int
var writer io.Writer

func startElement(n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		var attrFiled string
		for _, attr := range n.Attr {
			value := attr.Val
			if attr.Key == "href" {
				value = "'" + value + "'"
			}
			attrFiled += " " + attr.Key + "=" + value
		}
		var tail string
		if n.FirstChild == nil {
			tail = "/"
		}
		//fmt.Printf("%s %d\n",n.Data,depth)
		_, _ = fmt.Fprintf(writer, "%*s<%s%s%s>\n", depth<<1, "", n.Data, attrFiled, tail)
		depth++
		if n.FirstChild == nil {
			depth--
		}
	case html.TextNode:
		_, _ = fmt.Fprintf(writer, "%*s%s\n", depth<<1, "", n.Data)
	case html.CommentNode:
		_, _ = fmt.Fprintf(writer, "//%*s%s\n", depth<<1, "", n.Data)
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode && n.FirstChild != nil {
		depth--
		_, _ = fmt.Fprintf(writer, "%*s</%s>\n", depth<<1, "", n.Data)
	}
}
