package main

import "testing"

func BenchmarkPopCount0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount0(testInput)
	}
}

func BenchmarkPopCount1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount1(testInput)
	}
}

func BenchmarkPopCount2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount2(testInput)
	}
}

func BenchmarkPopCount3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount3(testInput)
	}
}
