package main

import (
	"bufio"
	"bytes"
	"fmt"
)

func main() {
	text := `HEAD
hello how are you
I'm fine, thank you.
In God WE TRUST.'`

	var wc WordCounter
	var lc LineCounter
	_, _ = fmt.Fprint(&wc, text)
	fmt.Println(wc)
	_, _ = fmt.Fprint(&lc, text)
	fmt.Println(lc)
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	cnt := 0
	for scanner.Scan() {
		cnt++
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	*c += WordCounter(cnt)
	return int(*c), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	cnt := 0
	for scanner.Scan() {
		cnt++
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	*c += LineCounter(cnt)
	return int(*c), nil
}
