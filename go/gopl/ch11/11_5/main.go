package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(split("Hello, world.", ","))
	fmt.Println(split("aab aab bba", "ab"))
	fmt.Println(split("すぐに返事ができなくてごめんなさいっ", "な"))
}

func split(s, sep string) (ret []string) {
	prev := 0
	for i := 0; i < len(s); {
		if strings.HasPrefix(s[i:], sep) {
			ret = append(ret, s[prev:i])
			prev = i + len(sep)
			i = prev
		} else {
			i++
		}
	}
	ret = append(ret, s[prev:])
	return
}
