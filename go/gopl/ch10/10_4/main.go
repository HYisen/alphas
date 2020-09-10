package main

import (
	"alphas/go/gopl/utility"
	"flag"
	"fmt"
	"strings"
)

var target = flag.String("target", "utility", "dependent package looking for")

func getAllPkg() []string {
	args := strings.Split("go list ...", " ")
	lines := utility.ExecAndGetStdOut(args...)
	split := strings.Split(lines, "\n")
	return split[0 : len(split)-1]
}

func getDeps(source string) []string {
	args := strings.Split(fmt.Sprintf("go list -f '{{.Deps}}' %s", source), " ")
	lines := utility.ExecAndGetStdOut(args...)
	lines = lines[2 : len(lines)-3]
	return strings.Split(lines, " ")
}

func main() {
	flag.Parse()

	for _, pkg := range getAllPkg() {
		//fmt.Printf("%s\n", pkg)

		for _, dep := range getDeps(pkg) {
			if strings.Contains(dep, *target) {
				fmt.Printf("%s -> %s\n", pkg, dep)
			}
			//fmt.Printf("\t%s\n", dep)
		}
	}

}
