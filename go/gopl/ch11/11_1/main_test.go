package main

import (
	"strings"
	"testing"
)

func TestCore(t *testing.T) {
	input := strings.NewReader("akka\n")
	var out strings.Builder
	var err strings.Builder
	core(input, &out, &err)
	const expect = `rune	count
'a'	2
'k'	2
'\n'	1

len	count
1	5
2	0
3	0
4	0
`
	if out.String() != expect {
		t.Errorf("expect\n%s\nactual\n%s\n", expect, out.String())
	}
	//fmt.Println(out.String())
}
