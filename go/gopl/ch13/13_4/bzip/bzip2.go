package bzip

import "C"

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
)

type writer struct {
	w io.Writer // underlying output stream

	lock sync.Mutex
}

func NewWriter(out io.Writer) io.WriteCloser {
	w := &writer{w: out}
	return w
}

func (w *writer) Write(data []byte) (int, error) {
	w.lock.Lock()
	defer func() { w.lock.Unlock() }()

	cmd := exec.Command("bzip2", "--stdout", "--compress")
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = w.w
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	err = cmd.Wait()
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (w *writer) Close() error {
	return nil
}
