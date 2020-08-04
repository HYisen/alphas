package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c0 := sha256.Sum256([]byte("ECHO"))
	c1 := sha256.Sum256([]byte("ECHo"))
	fmt.Printf("%x\n%x\ndiff=%d\n", c0, c1, CompareHash(c0, c1))
}

func CompareHash(c0, c1 [32]byte) int {
	cnt := 0
	for i := 0; i < 32; i++ {
		cnt += int(compare(c0[i], c1[i]))
	}
	return cnt
}

var db [256][256]uint8 // i > j

func init() {
	for i := range db {
		for j := 0; j < i; j++ {
			db[i][j] = calc(uint8(i), uint8(j))
		}
	}
}

// Will the compiler optimize it with registers?
func calc(l, r byte) uint8 {
	var cnt uint8
	for i := 0; i < 8; i++ {
		if l&r&1 == 1 {
			cnt++
		}
		l >>= 1
		r >>= 1
	}
	return cnt
}

func compare(b0, b1 byte) uint8 {
	if b0 == b1 {
		return 8
	}
	if b0 > b1 {
		return db[b0][b1]
	}
	return db[b1][b0]
}
