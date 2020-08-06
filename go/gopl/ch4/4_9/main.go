package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

func main() {
	data := map[string]int{} // use lower case word as key

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("failed to read : %v", err)
		}
		key := lower(scanner.Text())
		data[key]++
	}

	for word, cnt := range data {
		fmt.Printf("%8s=%4d\n", word, cnt)
	}
}

func lower(src string) string {
	runes := []rune(src)
	for i, r := range runes {
		runes[i] = unicode.ToLower(r)
	}
	return string(runes)
}
