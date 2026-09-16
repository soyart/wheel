package list

import (
	"cmp"
	"container/heap"
	"math/big"
	"reflect"
	"testing"

	"github.com/soyart/wheel"
)

const (
	minHeap = wheel.Ascending
	maxHeap = wheel.Descending
)

type foo[T cmp.Ordered] struct {
	name  string
	value T
}

// Implements data.Valuer[T]
func (f foo[T]) GetValue() T {
	return f.value
}

func TestPq(t *testing.T) {
	highest := foo[int]{name: "b", value: 100}
	lowest := foo[int]{name: "d", value: 0}

	items := []foo[int]{
		{name: "a", value: 69},
		highest,
		{name: "c", value: 12},
		lowest,
	}

	minHeapResult := testPop(t, minHeap, items)
	if minHeapResult != lowest {
		t.Fatalf("unexpected MinHeap results - expected %+v, got %+v\n", lowest, minHeapResult)
	}

	maxHeapResult := testPop(t, maxHeap, items)
	if maxHeapResult != highest {
		t.Fatalf("unexpected MaxHeap results - expected %+v, got %+v\n", highest, maxHeapResult)
	}

	testArbitaryUpdate(t)
}

// testPop orders foo[T] by its value field, since foo[T] itself is not cmp.Ordered.
func testPop[T cmp.Ordered](t *testing.T, order wheel.SortOrder, items []foo[T]) foo[T] {
	pq := NewPriorityQueueCustom(order, wheel.LessFuncBy(order, foo[T].GetValue))
	for _, item := range items {
		heap.Push(pq, item)
	}

	p := heap.Pop(pq)
	popped, ok := p.(foo[T])
	if !ok {
		t.Fatalf("type assertion to foo failed, got type: %s", reflect.TypeOf(p))
	}

	return popped
}

func testArbitaryUpdate(t *testing.T) {
	// Test with Valuer[float64]
	hundred := foo[float64]{name: "hundred", value: 100}
	seventy := foo[float64]{name: "seventy", value: 70}
	zero := foo[float64]{name: "zero", value: 0}
	foosFloat := []foo[float64]{
		{name: "a", value: 69},
		hundred,
		{name: "b", value: 71},
		zero,
		seventy,
	}

	pq := NewPriorityQueueCustom(maxHeap, wheel.LessFuncBy(maxHeap, foo[float64].GetValue))

	for _, item := range foosFloat {
		heap.Push(pq, item)
	}

	p := heap.Pop(pq)
	popped, ok := p.(foo[float64])
	if !ok {
		t.Fatalf("type assertion to foo failed, got type: %s", reflect.TypeOf(p))
	}
	if popped != hundred {
		t.Fatalf("unexpected MaxHeap results - expected %+v, got %+v\n", hundred, popped)
	}

	pq.ChangeOrdering(wheel.LessFuncBy(minHeap, foo[float64].GetValue))
	heap.Init(pq)
	p = heap.Pop(pq)
	popped, ok = p.(foo[float64])
	if !ok {
		t.Fatalf("type assertion to foo failed, got type: %s", reflect.TypeOf(p))
	}
	if popped != zero {
		t.Fatalf("unexpected MinHeap results - expected %+v, got %+v\n", zero, popped)
	}
}

type bar struct {
	val *big.Int
}

func (b *bar) GetValue() *big.Int {
	return b.val
}

// barLessFunc orders *bar by its *big.Int value via Cmp, since *bar is not wheel.CmpOrdered itself.
func barLessFunc(order wheel.SortOrder) wheel.LessFunc[*bar] {
	if order == wheel.Ascending {
		return func(items []*bar, i, j int) bool {
			return items[i].GetValue().Cmp(items[j].GetValue()) < 0
		}
	}
	return func(items []*bar, i, j int) bool {
		return items[i].GetValue().Cmp(items[j].GetValue()) > 0
	}
}

func TestPQCmp(t *testing.T) {
	a := &bar{val: big.NewInt(69)}
	b := &bar{val: big.NewInt(70)}
	c := &bar{val: big.NewInt(100)}
	d := &bar{val: big.NewInt(1000000)}

	t.Run("MaxHeap with Cmp", func(t *testing.T) {
		testPqCmpMax(t, []*bar{a, d, c, b}, d)
	})
	t.Run("MinHeap with Cmp", func(t *testing.T) {
		testPqCmpMin(t, []*bar{a, d, c, b}, a)
	})
}

func testPqCmpMax(t *testing.T, messy []*bar, max *bar) {
	maxPq := NewPriorityQueueCustom(maxHeap, barLessFunc(maxHeap))

	for _, item := range messy {
		heap.Push(maxPq, item)
	}
	if popped := heap.Pop(maxPq); popped != nil {
		actual := popped.(*bar)
		if actual != max {
			t.Fatalf("unexpected max heap result: expected %v, got %v\n", max.GetValue(), actual.GetValue())
		}
	}
}

func testPqCmpMin(t *testing.T, messy []*bar, min *bar) {
	minPq := NewPriorityQueueCustom(minHeap, barLessFunc(minHeap))

	for _, item := range messy {
		heap.Push(minPq, item)
	}
	if popped := heap.Pop(minPq); popped != nil {
		actual := popped.(*bar)
		if actual != min {
			t.Fatalf("unexpected min heap result: expected %v, got %v\n", min.GetValue(), actual.GetValue())
		}
	}
}
