package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	var array []string
	for k := 0; k < 1000; k++ {
		array = append(array, fmt.Sprintf("%4d", k))
	}

	test("one", func() {
		one(array...)
	})

	test("two", func() {
		two(array...)
	})

	test("one", func() {
		one(array...)
	})

	test("two", func() {
		two(array...)
	})
}

func test(name string, procedure func()) {
	start := time.Now()
	for k := 0; k < 1000; k++ {
		procedure()
	}
	fmt.Printf("%s cost %d ms\n", name, time.Since(start).Milliseconds())
}

func one(limb ...string) string {
	return strings.Join(limb, " ")
}

func two(limb ...string) string {
	var ret string
	for i, s := range limb {
		if i != 0 {
			ret += " "
		}
		ret = ret + s
	}
	return ret
}
