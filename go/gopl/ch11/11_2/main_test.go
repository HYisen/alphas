package main

import (
	"testing"
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

	oSet[14] = true
	if !nSet.Has(2) {
		t.Error("remove")
	}
}
