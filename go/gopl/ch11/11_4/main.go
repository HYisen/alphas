package main

import (
	"fmt"
	"unicode"
)

func main() {
	fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
}

func IsPalindrome(orig string) bool {
	var str []rune
	for _, r := range []rune(orig) {
		// remove punctuations & space
		if !unicode.IsPrint(r) || unicode.IsPunct(r) || unicode.IsSpace(r) {
			continue
		}

		str = append(str, unicode.ToLower(r))
	}

	if len(str) <= 1 {
		return true
	}

	head, tail := 0, len(str)-1
	for head < tail {
		// False positive of GoLand check, str is prevented from being nil by previous len check. issue/GO-9652
		if //goland:noinspection GoNilness
		str[head] != str[tail] {
			return false
		}
		head++
		tail--
	}
	return true
}
