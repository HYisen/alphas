package main

import (
	"alphas/go/gopl/utility"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

var targetID = "sec-origin-goals"
var targetClass = "copyright"

func main() {
	var result []string
	var sb strings.Builder

	dec := xml.NewDecoder(os.Stdin)
	var stack []string
	var matched bool
	var innerCount int

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local)

			// handle inner
			if matched {
				innerCount++
			}
			// handle gate
			if !matched && isMatch(tok) {
				//fmt.Printf("match at %s: %s\n", strings.Join(stack, " "), tok)
				matched = true
				sb.Reset()
			}
			// handle print
			if matched {
				sb.WriteString(fmt.Sprintf("<%s", tok.Name.Local))
				for _, attr := range tok.Attr {
					sb.WriteString(fmt.Sprintf(" %s=%s", attr.Name.Local, attr.Value))
				}
				sb.WriteString(">")
				//fmt.Println("add header "+sb.String())
				//fmt.Println(len(sb.String()))
			}
		case xml.EndElement:
			stack = stack[:len(stack)-1]

			// handle print
			if matched {
				sb.WriteString(fmt.Sprintf("</%s>", tok.Name.Local))
			}
			// handle gate
			if matched && innerCount == 0 {
				matched = false
				//fmt.Printf("add %d\n", len(sb.String()))
				//fmt.Println(sb.String())
				result = append(result, sb.String())
			}
			// handle inner
			if matched {
				innerCount--
			}
		case xml.CharData:
			if utility.ContainsAll(stack, os.Args[1:]) {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}

			if matched {
				sb.Write(tok)
			}
		}
	}

	for _, line := range result {
		fmt.Println(line)
	}
}

// isMatch check whether id or class is match in attr
func isMatch(e xml.StartElement) bool {
	for _, attr := range e.Attr {
		switch attr.Name.Local {
		case "id":
			if attr.Value == targetID {
				return true
			}
		case "class":
			if attr.Value == targetClass {
				return true
			}
		}
	}
	return false
}
