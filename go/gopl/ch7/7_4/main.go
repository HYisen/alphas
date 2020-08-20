package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
)

func NewReader(s string) *Reader {
	r := Reader{s}
	return &r
}

type Reader struct {
	s string
}

func (r Reader) Read(p []byte) (n int, err error) {
	fmt.Println("read")
	copy(p, r.s)
	return len(r.s), io.EOF
}

func main() {
	node, err := html.Parse(NewReader("<h>hello</h>"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(node.Data)
	fmt.Println(node.FirstChild.Data)
	fmt.Println(node.FirstChild.FirstChild.Data)
	fmt.Println(node.FirstChild.FirstChild.NextSibling.Data)
	fmt.Println(node.FirstChild.FirstChild.NextSibling.FirstChild.Data)
	fmt.Println(node.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.Data)
}
