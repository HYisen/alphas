package main

import (
	"fmt"
	"strings"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) == 1
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
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
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				limbs = append(limbs, fmt.Sprintf("%d", 64*i+j))
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
	word, bit := x/64, uint(x%64)
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
}
