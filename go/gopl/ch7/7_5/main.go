package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func LimitReader(r io.Reader, n int64) io.Reader {
	return &reader{
		orig: r,
		left: n,
	}
}

type reader struct {
	orig io.Reader
	left int64
}

func (r reader) Read(p []byte) (n int, err error) {
	if int64(len(p)) < r.left {
		read, err := r.orig.Read(p)
		r.left -= int64(read)
		return read, err
	}

	// As that in the io.Reader's document,
	// I can use all of p as scratch space during the call.
	_, _ = r.orig.Read(p)
	cnt := r.left
	r.left = 0
	return int(cnt), io.EOF
}

func main() {
	buff := new(bytes.Buffer)
	n, err := buff.ReadFrom(LimitReader(strings.NewReader("012345"), 2))
	fmt.Println(n, err, buff)
}
