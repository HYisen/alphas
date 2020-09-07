package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	head := make(chan int)
	prev := head
	var size int = 1e6
	for i := 0; i < size; i++ {
		neo := make(chan int)
		go action(prev, neo)
		prev = neo
	}

	fmt.Printf("init cost %v\n", time.Since(start))
	start = time.Now()

	head <- 0
	fmt.Println(<-prev)
	fmt.Printf("pass cost %v\n", time.Since(start))
}

func action(head <-chan int, tail chan<- int) {
	num := <-head
	//fmt.Println(num)
	num++
	tail <- num
}
