package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(comma("12345678"))
}

func comma(s string) string {
	var buf bytes.Buffer

	length := len(s)
	first := length % 3
	size := length / 3
	buf.WriteString(s[:first])
	for i := 0; i < size; i++ {
		buf.WriteRune(',')
		buf.WriteString(s[first : first+3])
		first += 3
	}

	return buf.String()
}
