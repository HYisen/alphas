package main

import (
	"fmt"
	"sync"
)

var once sync.Once
var pc [256]byte

const testInput = 19700101

func prepare() {
	fmt.Println("prepare")
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	once.Do(prepare)
	return int(pc[byte(x>>(0*8))]) +
		int(pc[byte(x>>(1*8))]) +
		int(pc[byte(x>>(2*8))]) +
		int(pc[byte(x>>(3*8))]) +
		int(pc[byte(x>>(4*8))]) +
		int(pc[byte(x>>(5*8))]) +
		int(pc[byte(x>>(6*8))]) +
		int(pc[byte(x>>(7*8))])
}

func main() {
	fmt.Println("one")
	fmt.Println(PopCount(testInput))
}
