package list

// Queue can be used in many scenarios, including when implementing a graph.
// See queueForGraph below for examples.

import (
	"testing"
)

func TestQueue(t *testing.T) {
	values0 := []uint8{0, 1, 200, 20}
	q0 := NewQueue[uint8]()
	testQueue(t, values0, q0)

	values1 := []string{"foo", "bar", "baz"}
	q1 := NewQueue[string]()
	testQueue(t, values1, q1)

	// Composite type queue - any comparable types should be ok in tests
	valuesComposite := []interface{}{1, 2, "last"}
	qComposite := NewQueue[interface{}]()
	// Test Push for composite queue
	qComposite.PushSlice(valuesComposite)
	// Test Pop for composite queue
	for i, qLen := 0, qComposite.Len(); i < qLen; i++ {
		expected := valuesComposite[i]
		p := qComposite.Pop()
		if p == nil {
			t.Fatalf("Queue.Pop failed - expected %v, got nil", expected)
		}
		value := *p
		if value != expected {
			t.Fatalf("Queue.Pop failed - expected %v, got %v", expected, value)
		}
	}
	// Test Queue.IsEmpty for composite queue
	if !qComposite.IsEmpty() {
		t.Fatalf("Queue.IsEmpty failed - expected to be emptied")
	}
}

func testQueue[T comparable](t *testing.T, values []T, q BasicList[T]) {
	// Test Push
	for _, expected := range values {
		q.Push(expected)
		v := q.Pop()
		if *v != expected {
			t.Fatalf("Queue.Pop failed - expected %v, got %v", expected, v)
		}
	}

	// Test Len
	for i, v := range values {
		q.Push(v)

		if newLen := q.Len(); newLen != i+1 {
			t.Fatalf("Queue.Len failed - expected %d, got %d", newLen, i+1)
		}
	}

	// Test Pop
	for i, qLen := 0, q.Len(); i < qLen; i++ {
		popped := q.Pop()
		expected := values[i]
		if *popped != expected {
			t.Fatalf("Queue.Pop failed - expected %v, got %v", expected, popped)
		}
	}

	// Test IsEmpty
	if !q.IsEmpty() {
		t.Fatal("Stack.IsEmpty failed - expected true")
	}

	// Test Pop after emptied
	v := q.Pop()
	t.Logf("value of Pop() after emptied: %v\n", v)
}
