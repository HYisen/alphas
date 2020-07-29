package main

import "fmt"

func main() {
	fmt.Println(isAnagrams("abc", "bca"))
	fmt.Println(isAnagrams("abc", "bcad"))
	fmt.Println(isAnagrams("a gentleman", "elegant man"))
	fmt.Println(isAnagrams("a gentlemans", "elegant man"))
	fmt.Println(isAnagrams("日往则月来", "月往则日来"))
	fmt.Println(isAnagrams("日往则月", "往则日来"))
}

func isAnagrams(l, r string) bool {
	data := make(map[rune]int)

	for _, ch := range l {
		data[ch]++
	}

	fmt.Println(data)

	for _, ch := range r {
		data[ch]--
	}

	fmt.Println(data)

	for _, cnt := range data {
		if cnt != 0 {
			return false
		}
	}
	return true
}
