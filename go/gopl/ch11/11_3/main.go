package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	//data := strings.Split("abba", "")
	data := strings.Split("1aba", "")

	fmt.Println(IsPalindrome(sort.StringSlice(data)))
}

func IsPalindrome(s sort.Interface) bool {
	if s.Len() <= 1 {
		return true
	}

	begin, end := 0, s.Len()-1
	for begin < end {
		if s.Less(begin, end) || s.Less(end, begin) {
			return false
		}
		begin++
		end--
	}
	return true
}
