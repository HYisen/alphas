package main

import (
	"alphas/go/gopl/ch2/2_1/tempconv"
	"fmt"
)

func main() {
	fmt.Printf("Brrrr! %v\n", tempconv.BoilingC)
	fmt.Printf("Brrrr! %v\n", tempconv.CToF(tempconv.BoilingC))
	fmt.Printf("Brrrr! %v\n", tempconv.CToK(tempconv.BoilingC))
}
