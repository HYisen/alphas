package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("fuck")
	fmt.Println(strings.Join(os.Args, " "))
}
