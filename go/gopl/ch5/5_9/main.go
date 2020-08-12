package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(expand("alex $alex use $GOD god", strings.ToUpper))
}

func expand(s string, f func(string) string) string {
	var data []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if word[0] == '$' {
			data = append(data, f(word[1:]))
		} else {
			data = append(data, word)
		}
	}
	return strings.Join(data, " ")
}
