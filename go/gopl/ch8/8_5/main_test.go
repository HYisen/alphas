package main

import "testing"

func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 0, -1)
	}
}

func Benchmark2K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 1, -1)
	}
}

func Benchmark2Kx2K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 2, -1)
	}
}

func Benchmark1T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 1)
	}
}

func Benchmark2T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 2)
	}
}

func Benchmark3T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 3)
	}
}

func Benchmark4T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 4)
	}
}

func Benchmark5T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 5)
	}
}

func Benchmark6T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 6)
	}
}

func Benchmark8T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 8)
	}
}

func Benchmark10T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 10)
	}
}
func Benchmark12T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 12)
	}
}
func Benchmark16T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 16)
	}
}
func Benchmark32T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 32)
	}
}
func Benchmark64T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 64)
	}
}
func Benchmark128T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 128)
	}
}
func Benchmark256T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 256)
	}
}
func Benchmark512T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 512)
	}
}
func Benchmark1024T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 1024)
	}
}
func Benchmark2048T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 2048)
	}
}
func Benchmark4096T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 4096)
	}
}
func Benchmark8192T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exec(false, 3, 8192)
	}
}
