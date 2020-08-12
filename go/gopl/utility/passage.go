package utility

import (
	"bufio"
	"strings"
)

func CountWords(text string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	cnt := 0
	for scanner.Scan() {
		cnt++
	}
	return cnt, scanner.Err()
}
