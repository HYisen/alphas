package utility

import (
	"testing"
)

func TestStack(t *testing.T) {
	var stack Stack
	stack.Push(1)
	if stack.Top().(int) != 1 {
		t.Fatalf("can not get newly pushed item")
	}
	stack.Push(2)
	if stack.Top().(int) != 2 {
		t.Fatalf("can not get newly pushed item")
	}
	stack.Pop()
	if stack.Top().(int) != 1 {
		t.Fatalf("can not pop correctly")
	}
	stack.Push(3)
	stack.Push(4)
	if stack.Top().(int) != 4 {
		t.Fatalf("can not push again")
	}
	stack.Pop()
	stack.Pop()
	if stack.Top().(int) != 1 {
		t.Fatalf("can not do difficult task")
	}
	stack.Pop()
}
