// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 101.

// Package treesort provides insertion sort using an unbalanced binary tree.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

//!+
type tree struct {
	value       int
	left, right *tree
}

func (t tree) String() string {
	var sb strings.Builder
	sb.WriteString("tree view:\n")
	rest := []*tree{&t}
	var nextDepthVacuum bool
	for !nextDepthVacuum {
		nextDepthVacuum = true
		var parents []string
		var children []*tree
		for i := 0; i < len(rest); i++ {
			one := rest[i]
			if one == nil {
				parents = append(parents, "X")
				children = append(children, nil, nil)
				continue
			}
			parents = append(parents, strconv.Itoa(one.value))
			children = append(children, one.left, one.right)
			if nextDepthVacuum && (one.left != nil || one.right != nil) {
				nextDepthVacuum = false
			}
		}
		rest = children
		sb.WriteString(strings.Join(parents, "|") + "\n")
	}
	return sb.String()
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
		fmt.Printf("after add %d %v", v, root)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

//!-

func main() {
	Sort([]int{1, 4, 6, 2, 8, 3})
}
