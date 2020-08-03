package main

import "fmt"

func reverse(ptr *[4]int) {
	for i, j := 0, len(*ptr)-1; i < j; i, j = i+1, j-1 {
		(*ptr)[i], (*ptr)[j] = (*ptr)[j], (*ptr)[i]
	}
}

func main() {
	arr := [...]int{1, 2, 4, 8}
	fmt.Println(arr)
	reverse(&arr)
	fmt.Println(arr)
}
