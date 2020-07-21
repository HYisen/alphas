package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("fuck")
	for i, arg := range os.Args[1:] {
		fmt.Printf("%4d %s\n", i+1, arg)
	}
}
