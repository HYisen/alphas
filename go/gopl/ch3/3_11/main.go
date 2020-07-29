package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(comma("12345678"))
	fmt.Println(comma("-12345678"))
	fmt.Println(comma("-12345678.925"))
	fmt.Println(comma("12345678.4"))
	fmt.Println(comma("+12345678.4"))
	fmt.Println(comma("+1234567"))
}

func comma(s string) string {
	var buf bytes.Buffer

	fmt.Printf("%s = ", s)

	if s[0] == '-' || s[0] == '+' {
		buf.WriteByte(s[0])
		s = s[1:]
	}

	dotIndex := strings.IndexByte(s, '.')
	var suffix string
	if dotIndex != -1 {
		suffix = s[dotIndex:]
		s = s[:dotIndex]
	}

	length := len(s)
	first := length % 3
	size := length / 3
	buf.WriteString(s[:first])
	for i := 0; i < size; i++ {
		buf.WriteRune(',')
		buf.WriteString(s[first : first+3])
		first += 3
	}

	buf.WriteString(suffix)

	return buf.String()
}
