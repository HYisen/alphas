package main

import (
	"alphas/go/gopl/ch7/7_13A14/eval"
	"fmt"
	"log"
)

func main() {
	exec("1 + 2*4 - 5/2 + x*1000", eval.Env{"x": 3.14})
	exec("1+min(x,8,3.8,-5)", eval.Env{"x": 3.14})
}

func exec(s string, env eval.Env) {
	expr, err := eval.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(expr)
	fmt.Println(expr.Eval(env))
}
