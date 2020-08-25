package main

import (
	"alphas/go/gopl/ch7/7_13/eval"
	"fmt"
	"log"
)

func main() {
	expr, err := eval.Parse("1+2*4-5/2+x*1000")
	if err != nil {
		log.Fatal(err)
	}
	env := eval.Env{"x": 3.14}
	fmt.Println(expr)
	fmt.Println(expr.Eval(env))
}
