package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRoute(t *testing.T) {
	oSet := make(map[int]bool)
	var nSet IntSet

	nSet.Add(2)

	oSet[14] = true
	if !nSet.Has(2) {
		t.Error("add")
	}

	nSet.Remove(2)

	fmt.Println(&nSet)

	oSet[14] = true
	if nSet.Has(2) {
		t.Error("remove")
	}
}

func BenchmarkAdd(b *testing.B) {
	benchmarks := []struct {
		wordSize int
	}{
		{wordSize: 64},
		{wordSize: 32},
		{wordSize: 16},
		{wordSize: 24},
		{wordSize: 48},
	}

	input := genRandomInt16Array(1e6)

	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("wordSize%d", bm.wordSize), func(b *testing.B) {
			bitSize = bm.wordSize
			for i := 0; i < b.N; i++ {
				var set IntSet
				for _, num := range input {
					set.Add(num)
				}
			}
		})
	}

	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			set := map[int]bool{}
			for _, num := range input {
				set[num] = true
			}
		}
	})
}

func BenchmarkUnionWith(b *testing.B) {
	benchmarks := []struct {
		wordSize int
	}{
		{wordSize: 64},
		{wordSize: 32},
		{wordSize: 16},
		{wordSize: 24},
		{wordSize: 48},
	}

	var input []IntSet
	var nInput []map[int]bool
	for i := 0; i < 1000; i++ {
		array := genRandomInt16Array(1000)

		var set IntSet
		set.AddAll(array...)
		input = append(input, set)

		nSet := make(map[int]bool)
		for _, num := range array {
			nSet[num] = true
		}
		nInput = append(nInput, nSet)
	}

	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("wordSize%d", bm.wordSize), func(b *testing.B) {
			bitSize = bm.wordSize
			for i := 0; i < b.N; i++ {
				var set IntSet
				for _, inputSet := range input {
					set.UnionWith(&inputSet)
				}
			}
		})
	}

	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			set := map[int]bool{}
			for _, inputSet := range nInput {
				for k, v := range inputSet {
					if v {
						set[k] = true
					}
				}
			}
		}
	})
}

// Whoever care the difference between slice and array.
// With the help of constexpr and std::array, I can make it a fix length array, in C++.
func genRandomInt16Array(length int) []int {
	const maxExclusive = 1 << 16
	var ret []int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		ret = append(ret, r.Intn(maxExclusive))
	}
	return ret
}
