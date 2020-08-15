package main

import (
	"fmt"
	"strconv"
)

func main() {
	defer func() {
		r := recover()
		s := fmt.Sprint(r)
		fmt.Println(s)
		num, _ := strconv.ParseInt(s, 10, 32)
		fmt.Printf("=%d\n", num)
	}()

	square(4)
}

func square(num int) {
	panic(num * num)
}
