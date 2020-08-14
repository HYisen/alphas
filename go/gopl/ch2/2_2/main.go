package main

import (
	"alphas/go/gopl/ch2/2_1/tempconv"
	"alphas/go/gopl/ch2/2_2/distconv"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	var input []float64

	switch len(args) {
	case 1:
		input = getInputFromCin()
	default:
		input = getInputFromArg(args)
	}

	for _, val := range input {
		fmt.Printf("\nfor input %g:\n", val)

		c := tempconv.Celsius(val)
		f := tempconv.Fahrenheit(val)
		fmt.Printf("%v=%v\t%v=%v\n", c, tempconv.CToF(c), f, tempconv.FToC(f))

		ft := distconv.Feet(val)
		m := distconv.Meter(val)
		fmt.Printf("%v=%v\t%v=%v\n", ft, distconv.FToM(ft), m, distconv.MToF(m))
	}
}

func extract(src []string) []float64 {
	var ret []float64

	for _, str := range src {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "bad input "+str)
			os.Exit(1)
		}
		ret = append(ret, val)
	}

	return ret
}

func getInputFromArg(args []string) []float64 {
	return extract(args[1:])
}

func getInputFromCin() []float64 {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()

	var ret []string

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return extract(ret)
}
