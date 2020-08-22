package tempconv

import (
	"flag"
	"fmt"
)

type celsiusFlag struct {
	Celsius
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	_, err := fmt.Sscanf(s, "%f%s", &value, &unit)
	if err != nil {
		return fmt.Errorf("can not parse %s:%v", s, err)
	}

	switch unit {
	case "C", "℃":
		f.Celsius = Celsius(value)
	case "F", "℉":
		f.Celsius = FToC(Fahrenheit(value))
	case "K", "K":
		f.Celsius = KToC(Kelvin(value))
	default:
		return fmt.Errorf("bad unit %s on value %f", unit, value)
	}
	return nil
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
