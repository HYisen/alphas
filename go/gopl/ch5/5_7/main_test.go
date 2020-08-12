package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "shall print href", input: `
<html><head/><body><a href='z.cn'/></body></html>`,
			expect: `<html>
  <head/>
  <body>
    <a href='z.cn'/>
  </body>
</html>
`}, {name: "shall handle text node", input: `
<html><head/><body><p>fuck</p></body></html>`,
			expect: `<html>
  <head/>
  <body>
    <p>
      fuck
    </p>
  </body>
</html>
`}, {name: "shall shrink element without child", input: `
<html><head/><body></body></html>`,
			expect: `<html>
  <head/>
  <body/>
</html>
`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			doc, err := html.Parse(strings.NewReader(test.input))
			if err != nil {
				err = fmt.Errorf("parsing HTML: %s", err)
				t.Fatal(err)
			}

			var sb strings.Builder
			writer = &sb

			forEachNode(doc, startElement, endElement)

			if sb.String() != test.expect {
				t.Fatalf("\nexpect:\n%s\nactual:\n%s\n", test.expect, sb.String())
			}
		})
	}
}
