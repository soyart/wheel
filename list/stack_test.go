package list

import "testing"

func TestStack(t *testing.T) {
	values0 := []uint8{0, 1, 2, 3}
	stack0 := NewStack[uint8]()
	sstack0 := NewStack[uint8]()
	uintStacks := []BasicList[uint8]{stack0, sstack0}
	for _, s := range uintStacks {
		testStack(t, values0, s)
	}

	values1 := []string{"one"}
	stack1 := NewStack[string]()
	sstack1 := NewStack[string]()
	stringStacks := []BasicList[string]{stack1, sstack1}
	for _, s := range stringStacks {
		testStack(t, values1, s)
	}

	// Composite type stack - any comparable types should be ok in this tests
	valuesComposite := []interface{}{1, true, "second last"}
	stackComposite := NewStack[interface{}]()
	// Test Push for composite queue
	for _, value := range valuesComposite {
		stackComposite.Push(value)
	}
	// Test Pop for composite queue
	for i, qLen := 0, stackComposite.Len(); i < qLen; i++ {
		expected := valuesComposite[qLen-i-1]
		p := stackComposite.Pop()
		if p == nil {
			t.Fatalf("Queue.Pop failed - expected %v, got nil", expected)
		}
		value := *p
		if value != expected {
			t.Fatalf("Queue.Pop failed - expected %v, got %v", expected, value)
		}
	}
	// Test Queue.IsEmpty for composite queue
	if !stackComposite.IsEmpty() {
		t.Fatalf("Queue.IsEmpty failed - expected to be emptied")
	}
}

func testStack[T comparable](t *testing.T, values []T, stack BasicList[T]) {
	// Test Push
	for _, expected := range values {
		stack.Push(expected)
		v := stack.Pop()
		if v == nil {
			t.Fatalf("Stack.Pop failed - expect %v, got nil", expected)
		}
		if *v != expected {
			t.Fatalf("Stack.Pop failed - expected %v, got %v", expected, v)
		}
	}

	// Test Len
	for i, v := range values {
		stack.Push(v)

		if newLen := stack.Len(); newLen != i+1 {
			t.Fatalf("Stack.Len failed - expected %d, got %d", newLen, i+1)
		}
	}

	// Test Pop
	for i, qLen := 0, stack.Len(); i < qLen; i++ {
		popped := stack.Pop()
		expected := values[qLen-i-1]
		if *popped != expected {
			t.Fatalf("Stack.Pop failed - expected %v, got %v", expected, popped)
		}
	}

	// Test IsEmpty
	if !stack.IsEmpty() {
		t.Fatal("Stack.IsEmpty failed - expected true")
	}

	// Test Pop after emptied
	v := stack.Pop()
	t.Logf("value of Pop() after emptied: %v\n", v)
}
