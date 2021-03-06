// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bzip_test

import (
	"alphas/go/gopl/ch13/13_4/bzip"
	"bytes"
	"compress/bzip2" // reader
	"io"
	"sync"
	"testing"
)

func TestBzip2(t *testing.T) {
	var compressed, uncompressed bytes.Buffer
	w := bzip.NewWriter(&compressed)

	// Write a repetitive message in a million pieces,
	// compressing one copy but not the other.
	tee := io.MultiWriter(w, &uncompressed)
	for i := 0; i < 10; i++ {
		io.WriteString(tee, "hello")
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	// Check the size of the compressed stream.
	if got, want := compressed.Len(), 410; got != want {
		t.Errorf("10 hellos compressed to %d bytes, want %d", got, want)
	}

	// Decompress and compare with original.
	var decompressed bytes.Buffer
	io.Copy(&decompressed, bzip2.NewReader(&compressed))
	if !bytes.Equal(uncompressed.Bytes(), decompressed.Bytes()) {
		t.Error("decompression yielded a different message")
	}
}

const thread = 4

func TestBzip2Concurrent(t *testing.T) {
	var compressed, uncompressed bytes.Buffer
	w := bzip.NewWriter(&compressed)

	// Write a repetitive message in a million pieces,
	// compressing one copy but not the other.
	tee := io.MultiWriter(w, &uncompressed)

	var wg sync.WaitGroup
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				io.WriteString(tee, "hello")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}
