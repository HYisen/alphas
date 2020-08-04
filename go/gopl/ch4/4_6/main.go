package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	data := []byte("a a  a\t a\t\tb")
	fmt.Println(data)
	fmt.Println(string(data))
	data = Squash(data)
	fmt.Println(data)
	fmt.Println(string(data))
}

func Squash(src []byte) []byte {
	head := 0
	prevIsSpace := false
	for tail := 0; tail < len(src); {
		r, size := utf8.DecodeRune(src[tail:])
		if unicode.IsSpace(r) {
			if !prevIsSpace {
				src[head] = ' '
				head += 1
			}
			prevIsSpace = true
		} else {
			prevIsSpace = false
			copy(src[head:], src[tail:tail+size])
			head += size
		}
		tail += size
	}
	return src[:head]
}
