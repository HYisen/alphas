package main

import (
	"fmt"
)

var pc [256]byte

const testInput = 19700101

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount0(x uint64) int {
	return int(pc[byte(x>>(0*8))]) +
		int(pc[byte(x>>(1*8))]) +
		int(pc[byte(x>>(2*8))]) +
		int(pc[byte(x>>(3*8))]) +
		int(pc[byte(x>>(4*8))]) +
		int(pc[byte(x>>(5*8))]) +
		int(pc[byte(x>>(6*8))]) +
		int(pc[byte(x>>(7*8))])
}

func PopCount1(x uint64) int {
	cnt := 0
	for i := 0; i < 8; i++ {
		cnt += int(pc[byte(x>>uint(i*8))])
	}
	return cnt
}

func PopCount2(x uint64) int {
	cnt := 0
	for i := 0; i < 64; i++ {
		if (x>>uint(i))&1 == 1 {
			cnt += 1
		}
	}
	return cnt
}

func PopCount3(x uint64) int {
	cnt := 0
	for x != 0 {
		x &= x - 1
		cnt++
	}
	return cnt
}

func main() {
	fmt.Println(PopCount0(testInput))
	fmt.Println(PopCount1(testInput))
	fmt.Println(PopCount2(testInput))
	fmt.Println(PopCount3(testInput))
}
