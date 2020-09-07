package main

import (
	"runtime"
	"testing"
)

func Benchmark1T(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark2T(b *testing.B) {
	runtime.GOMAXPROCS(2)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark3T(b *testing.B) {
	runtime.GOMAXPROCS(3)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark4T(b *testing.B) {
	runtime.GOMAXPROCS(4)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark5T(b *testing.B) {
	runtime.GOMAXPROCS(5)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark6T(b *testing.B) {
	runtime.GOMAXPROCS(6)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark8T(b *testing.B) {
	runtime.GOMAXPROCS(8)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}

func Benchmark10T(b *testing.B) {
	runtime.GOMAXPROCS(10)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark12T(b *testing.B) {
	runtime.GOMAXPROCS(12)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark16T(b *testing.B) {
	runtime.GOMAXPROCS(16)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark32T(b *testing.B) {
	runtime.GOMAXPROCS(32)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark64T(b *testing.B) {
	runtime.GOMAXPROCS(64)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark128T(b *testing.B) {
	runtime.GOMAXPROCS(128)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark256T(b *testing.B) {
	runtime.GOMAXPROCS(256)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark512T(b *testing.B) {
	runtime.GOMAXPROCS(512)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
func Benchmark1024T(b *testing.B) {
	runtime.GOMAXPROCS(1024)
	for i := 0; i < b.N; i++ {
		exec(false)
	}
}
