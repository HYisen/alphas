package main

import (
	"alphas/go/gopl/utility"
	"fmt"
	"strings"
)

const bitSize = 16

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/bitSize, uint(x%bitSize)
	return word < len(s.words) && s.words[word]&(1<<bit) == 1
}

func (s *IntSet) Add(x int) {
	word, bit := x/bitSize, uint(x%bitSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tWord := range t.words {
		if i < len(s.words) {
			s.words[i] |= tWord
		} else {
			s.words = append(s.words, tWord)
		}
	}
}

func (s *IntSet) String() string {
	var limbs []string
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitSize; j++ {
			if word&(1<<uint(j)) != 0 {
				limbs = append(limbs, fmt.Sprintf("%d", bitSize*i+j))
			}
		}
	}
	return "{" + strings.Join(limbs, " ") + "}"
}

func (s *IntSet) Len() int {
	cnt := 0
	for _, word := range s.words {
		for word != 0 {
			word &= word - 1
			cnt++
		}
	}
	return cnt
}

func (s *IntSet) Remove(x int) {
	word, bit := x/bitSize, uint(x%bitSize)
	s.words[word] &= ^(1 << bit)
}

func (s *IntSet) Clear() {
	s.words = []uint64{}
}

func (s *IntSet) Copy() *IntSet {
	var neo IntSet
	for _, word := range s.words {
		neo.words = append(neo.words, word)

	}
	return &neo
}

func (s *IntSet) AddAll(x ...int) {
	for _, one := range x {
		s.Add(one)
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i := range s.words {
		if i >= len(t.words) {
			break
		}
		s.words[i] &= t.words[i]
	}
	for i, tWord := range t.words {
		if i < len(s.words) {
			s.words[i] |= tWord
		} else {
			s.words = append(s.words, tWord)
		}
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	tmp := s.Copy()
	tmp.SymmetricDifferenceWith(t)
	for i := range s.words {
		s.words[i] &= tmp.words[i]
	}
}

func (s *IntSet) SymmetricDifferenceWith(t *IntSet) {
	// for equivalent range word
	min, max := utility.MinAndMax(len(s.words), len(t.words))
	for i := 0; i < min; i++ {
		s.words[i] ^= t.words[i]
	}

	// for larger t word
	if len(s.words) < len(t.words) {
		for i := min; i < max; i++ {
			s.words = append(s.words, t.words[i])
		}
	}
}

func (s *IntSet) Elems() []int {
	var ret []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitSize; j++ {
			if word&(1<<uint(j)) != 0 {
				ret = append(ret, bitSize*i+j)
			}
		}
	}
	return ret
}

func main() {
	var a IntSet

	a.Add(1)
	a.Add(14)
	a.Add(32)
	fmt.Println(&a)
	fmt.Println(a.Len())

	a.Remove(8)
	a.Remove(14)
	fmt.Println(&a)
	fmt.Println(a.Len())

	fmt.Println("copy")
	pb := a.Copy()
	fmt.Println(&a)
	fmt.Println(pb)
	fmt.Println("add a")
	a.Add(55)
	fmt.Println(&a)
	fmt.Println(pb)

	fmt.Println("clear")
	a.Clear()
	fmt.Println(&a)
	fmt.Println(a.Len())
	fmt.Println(pb)

	fmt.Println("AddAll")
	a.AddAll(1, 2, 3, 4)
	fmt.Println(a.String())

	var one, two IntSet
	one.AddAll(1, 2, 3)
	two.AddAll(2, 3, 4)
	tmpA, tmpB := one.Copy(), two.Copy()
	tmpA.DifferenceWith(&two)
	tmpB.DifferenceWith(&one)
	fmt.Printf("Diff %v,%v\n", tmpA, tmpB)
	tmpA, tmpB = one.Copy(), two.Copy()
	tmpA.SymmetricDifferenceWith(&two)
	tmpB.SymmetricDifferenceWith(&one)
	fmt.Printf("SymDiff %v,%v\n", tmpA, tmpB)

	one.Clear()
	one.AddAll(18, 2, 3)
	for _, item := range one.Elems() {
		fmt.Println(item)
	}
}
