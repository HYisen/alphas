package main

import (
	"alphas/go/gopl/utility"
	"os"
)

func main() {
	node := utility.ParseXML(os.Stdin)
	//node := utility.ParseXML(strings.NewReader("<p>HW</p>"))
	utility.Print(0, node)
}
