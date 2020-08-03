package main

import "fmt"

func main() {
	data := []string{"apple", "apple", "apple", "banana", "cake", "duke", "duke", "enterprise"}
	fmt.Println(data)
	data = removeAdjacentDuplicates(data)
	fmt.Println(data)
}

func removeAdjacentDuplicates(src []string) []string {
	prev := 0
	for curr := 1; curr < len(src); curr++ {
		if src[curr] != src[prev] {
			prev++
			src[prev] = src[curr]
		}
	}
	return src[:prev+1]
}
