package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

func main() {
	data := make(map[string]int)

	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed to read : %v", err)
		}

		if unicode.IsLetter(r) {
			data["letter"]++
		}
		if unicode.IsNumber(r) {
			data["number"]++
		}
		if unicode.IsPunct(r) {
			data["punct"]++
		}
		if unicode.IsUpper(r) {
			data["upper"]++
		}
		if unicode.IsPrint(r) {
			data["print"]++
		}
	}

	for name, cnt := range data {
		fmt.Printf("%8s=%4d\n", name, cnt)
	}
}
