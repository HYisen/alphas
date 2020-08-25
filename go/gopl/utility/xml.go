package utility

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

type Node interface{} // CharData or *Element
type CharData string
type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

var spaceChains = []string{" "}

func getSpace(width int) string {
	if width < 0 {
		panic(fmt.Errorf("require space with minus depth %d", width))
	}
	if width == 0 {
		return ""
	}
	for width > len(spaceChains) {
		spaceChains = append(spaceChains, spaceChains[len(spaceChains)-1]+" ")
	}
	return spaceChains[width-1]
}

func printWithDepth(depth int, line string) {
	fmt.Println(getSpace(depth*2) + line)
}

func Print(depth int, node Node) {
	switch node := node.(type) {
	case CharData:
		printWithDepth(depth, string(node))
	case *Element:
		//fmt.Println("meet elem")
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("<%s", node.Type.Local))
		for _, attr := range node.Attr {
			sb.WriteString(fmt.Sprintf(" %s=%s", attr.Name.Local, attr.Value))
		}
		sb.WriteString(">")
		printWithDepth(depth, sb.String())
		for _, child := range node.Children {
			Print(depth+1, child)
		}
		printWithDepth(depth, fmt.Sprintf("</%s>", node.Type.Local))
	default:
		//forget this would cause catastrophic silent exception swallow during error handling
		panic("fuck type")
	}
}

func ParseXML(r io.Reader) Node {
	dec := xml.NewDecoder(r)
	// name fulfilled dummy node is helpful to debug
	dummy := Element{
		Type: xml.Name{
			Local: "dummy",
		},
		Attr:     nil,
		Children: []Node{},
	}
	var stack Stack
	stack.Push(&dummy)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			//fmt.Println("start")
			current := Element{
				Type:     tok.Name,
				Attr:     tok.Attr,
				Children: []Node{},
			}
			//fmt.Println(current.Type.Local)
			stack.Push(&current)
		case xml.EndElement:
			//fmt.Println("end")
			child := stack.Top().(*Element)
			//fmt.Println(child.Type.Local)
			stack.Pop()
			parent := stack.Top().(*Element)
			//fmt.Println(parent.Type.Local)
			parent.Children = append(parent.Children, child)
		case xml.CharData:
			//fmt.Println(string(tok))
			parent := stack.Top().(*Element)
			parent.Children = append(parent.Children, CharData(tok))
		}
	}

	return stack.Top().(*Element).Children[0]
}
