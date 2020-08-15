package main

import "fmt"

func main() {
	fmt.Println(max(1, 2, 3, 4, 5))
	fmt.Println(max(1, 5))
	fmt.Println(min(6, 5, 3, 8))
	fmt.Println(min(-5))
}

func max(first int, others ...int) int {
	ret := first
	for _, other := range others {
		if ret < other {
			ret = other
		}
	}
	return ret
}

func min(first int, others ...int) int {
	ret := first
	for _, other := range others {
		if ret > other {
			ret = other
		}
	}
	return ret
}
