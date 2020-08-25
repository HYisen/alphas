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

// ContainsAll reports whether array x contains every item in array y, in order.
func ContainsAll(x, y []string) bool {
	if len(x) < len(y) {
		return false
	}

	for i := range y {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}
