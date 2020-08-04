package main

import "fmt"

func main() {
	data := []int{5, 4, 3, 2, 1}
	fmt.Println(data)
	data = rotate(data[:], 3)
	fmt.Println(data)
}

func rotate(src []int, toLeftRange int) []int {
	return append(src[toLeftRange:], src[:toLeftRange]...)
}
