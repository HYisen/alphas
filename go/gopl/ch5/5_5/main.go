package main

import (
	"alphas/go/gopl/utility"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func main() {
	fmt.Println(CountWordsAndImages("https://www.taobao.com"))
}

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n.Type == html.ElementNode {
		if n.Data == "img" {
			images++
		}
	} else if n.Type == html.TextNode {
		neoWords, _ := utility.CountWords(n.Data)
		words += neoWords
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		neoWords, neoImages := countWordsAndImages(child)
		words += neoWords
		images += neoImages
	}
	return
}
