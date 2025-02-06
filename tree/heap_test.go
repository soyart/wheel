package tree

import (
	"math/big"
	"sort"
	"testing"

	"github.com/soyart/wheel"
)

func TestHeapifyUp(t *testing.T) {
	ints := []int{6, 1, 3, 3, 2, 4, 5}
	pq := NewHeap[int](wheel.Descending)

	for i := range ints {
		pq.Push(ints[i])
	}

	pq.Push(3)
	pq.Push(0)
	pq.Push(9)

	l := len(pq.Items)
	for curr := range pq.Items {
		// t.Log("index", curr, "value", pq.Items[curr].GetValue())

		parent := ParentIdx(curr)
		if parent >= 0 {
			if pq.LessFunc(pq.Items, curr, parent) {
				valueCurr := pq.Items[curr].GetValue()
				valueParent := pq.Items[parent].GetValue()

				t.Logf(
					"node at %d is less than parent at %d: node %v vs parent %v",
					curr, parent, valueCurr, valueParent,
				)

				t.Fatalf("Unexpected value")
			}
		}

		childLeft := LeftChildIdx(curr)
		if childLeft >= l {
			continue
		}

		if !pq.LessFunc(pq.Items, curr, childLeft) {
			valueCurr := pq.Items[curr].GetValue()
			valueLeft := pq.Items[childLeft].GetValue()

			t.Logf(
				"node at %d is less than left child at %d: node %v vs parent %v",
				curr, childLeft, valueCurr, valueLeft,
			)

			t.Fatalf("Unexpected value")
		}

		childRight := childLeft + 1
		if childRight >= l {
			continue
		}

		if !pq.LessFunc(pq.Items, curr, childRight) {
			valueCurr := pq.Items[curr].GetValue()
			valueRight := pq.Items[childRight].GetValue()

			t.Logf(
				"node at %d is less than left child at %d: node %v vs parent %v",
				curr, childRight, valueCurr, valueRight,
			)

			t.Fatalf("Unexpected heap shape")
		}
	}
}

func TestHeapifyDown(t *testing.T) {
	ints := []int{6, 1, 3, 3, 2, 4, 5}
	pq := NewHeap[int](wheel.Descending)

	for i := range ints {
		pq.Push(ints[i])
	}

	intsHeap := make([]int, len(pq.Items))
	for i := range pq.Items {
		intsHeap[i] = pq.Items[i].GetValue()
	}

	// t.Log("intsHeap", intsHeap)

	var c int
	for !pq.IsEmpty() {
		max := wheel.MaxValuer(pq.Items)
		popped := pq.PopValue()

		items := make([]int, pq.Len())
		for i := range pq.Items {
			items[i] = pq.Items[i].GetValue()
		}

		// t.Log("l", l, "max", max, "pop", popped, "items", items)

		if max != popped {
			t.Logf("pop #%d expecting max %d, got %d", c+1, max, popped)
			t.Fatal("Unexpected popped value")
		}

		c++
	}
}

func TestHeapCmp(t *testing.T) {
	ints := []int64{2, 1, 4, 0, 6, 70, 1, 6}
	min := wheel.Min(ints...)

	intsBig := make([]*big.Int, len(ints))
	for i := range ints {
		intsBig[i] = big.NewInt(ints[i])
	}

	h := NewHeapCmp[*big.Int](wheel.Ascending)
	for i := range intsBig {
		h.Push(intsBig[i])
	}

	zero := h.PeekValue()
	if v := zero.Int64(); v != min {
		t.Fatalf("unexpected root node, expecting %d, got %d", min, v)
	}

	sort.Slice(ints, func(i, j int) bool {
		return ints[i] < ints[j]
	})

	c := 0
	for !h.IsEmpty() {
		popped := h.PopValue().Int64()
		expected := ints[c]

		if popped != expected {
			t.Logf("Pop #%d:unexpected value: expecting %d, got %d", c, expected, popped)
			t.Fatalf("unexpected value")
		}

		c++
	}
}
