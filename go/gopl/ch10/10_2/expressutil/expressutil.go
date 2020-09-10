package expressutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type MatchFunc func(data []byte) bool
type DepressFunc func(r io.Reader, w io.Writer) error

type item struct {
	matchFunc   MatchFunc
	depressFunc DepressFunc
}

var dict = make(map[string]item)

func RegisterFormat(name string, matchFunc MatchFunc, depressFunc DepressFunc) {
	if _, ok := dict[name]; ok {
		panic("existed register name " + name)
	}
	dict[name] = item{
		matchFunc:   matchFunc,
		depressFunc: depressFunc,
	}
}

func Depress(compressedData []byte) (string, error) {
	for name, item := range dict {
		if !item.matchFunc(compressedData) {
			continue
		}
		_, _ = fmt.Fprintln(os.Stderr, "match type "+name)

		reader := bytes.NewReader(compressedData)
		var sb strings.Builder
		err := item.depressFunc(reader, &sb)
		if err != nil {
			return "", err
		}
		return sb.String(), nil
	}
	return "", fmt.Errorf("unsupported compress type")
}
