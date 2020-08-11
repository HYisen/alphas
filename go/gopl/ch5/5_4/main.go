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
	visit(doc)
	fmt.Println(data)
}

var data = make(map[string][]string)

func add(key, val string) {
	data[key] = append(data[key], val)
}

var currentBlockName string

func visit(node *html.Node) {
	if node.Type == html.ElementNode {
		dealAttr(node, "a", "href")
		dealAttr(node, "img", "src")
		currentBlockName = node.Data
		if node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					add("img", attr.Val)
				}
			}
		}
	} else if node.Type == html.TextNode {
		if currentBlockName == "script" {
			add("script", node.Data)
		} else if currentBlockName == "style" {
			add("style", node.Data)
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		visit(child)
	}
}

func dealAttr(node *html.Node, blockName, attrName string) {
	if node.Data == blockName {
		for _, attr := range node.Attr {
			if attr.Key == attrName {
				add(blockName, attr.Val)
			}
		}
	}
}
