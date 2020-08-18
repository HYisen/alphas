package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	sb := strings.Builder{}
	neo, cnt := CountingWriter(&sb)
	_, _ = fmt.Fprint(neo, "放肆")
	fmt.Println(sb.String())
	fmt.Println(cnt)
	fmt.Println(*cnt)
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	ret := wrapper{
		inner: w,
		cnt:   0,
	}
	return &ret, &(ret.cnt)
}

type wrapper struct {
	inner io.Writer
	cnt   int64
}

func (w *wrapper) Write(p []byte) (int, error) {
	w.cnt += int64(len(p))
	return w.inner.Write(p)
}
