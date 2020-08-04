package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	data := []byte("イッツオーライ!ETC")
	fmt.Println(data)
	fmt.Println(string(data))
	reverse(data)
	fmt.Println(data)
	fmt.Println(string(data))
}

func reverse(src []byte) {
	for i := 0; i < len(src); {
		r, size := utf8.DecodeRune(src)
		copy(src, src[size:len(src)-i])
		copy(src[len(src)-i-size:], string(r))
		i += size
		//fmt.Println(i)
		//fmt.Println(src)
		//fmt.Println(string(src))
	}
}
