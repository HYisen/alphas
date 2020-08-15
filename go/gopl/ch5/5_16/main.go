package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(join("|", "a", "b", "c"))
	fmt.Println(join("|"))
	fmt.Println(join("&", "apple"))
}

func join(sep string, limbs ...string) string {
	var sb strings.Builder
	for i, limb := range limbs {
		if i != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(limb)
	}
	return sb.String()
}
