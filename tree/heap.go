package tree

import (
	"cmp"

	"github.com/soyart/wheel"
)

// Heap is a binary heap implementation backed by Go slice.
type Heap[T any] struct {
	Items    []T
	LessFunc wheel.LessFunc[T]
}

type HeapOption func(*HeapOptions)

type HeapOptions struct {
	preAlloc int
}

func HeapPreAlloc(size int) HeapOption {
	return func(opt *HeapOptions) {
		opt.preAlloc = size
	}
}

func NewHeap[T cmp.Ordered](order wheel.SortOrder, opts ...HeapOption) *Heap[T] {
	return NewHeapCustom[T](order, wheel.FactoryLessFuncOrdered[T](order), opts...)
}

func NewHeapCmp[T wheel.CmpOrdered[T]](order wheel.SortOrder, opts ...HeapOption) *Heap[T] {
	return NewHeapCustom[T](order, wheel.FactoryLessFuncCmp[T](order), opts...)
}

// NewHeapCustom builds a Heap[T] ordered by the given lessFunc, for T that
// isn't cmp.Ordered/CmpOrdered itself (e.g. a struct ordered by one of its
// fields). Pair it with wheel.LessFuncBy to order by an extracted key.
func NewHeapCustom[T any](order wheel.SortOrder, lessFunc wheel.LessFunc[T], opts ...HeapOption) *Heap[T] {
	options := parseOptions(opts...)
	return &Heap[T]{
		Items:    make([]T, 0, options.preAlloc),
		LessFunc: lessFunc,
	}
}

func NewHeapFrom[T cmp.Ordered](order wheel.SortOrder, items []T, opts ...HeapOption) *Heap[T] {
	// Default to pre-allocating slice of size len(items)
	// If given options define larger size, use that size
	options := parseOptions(opts...)
	preAlloc := HeapPreAlloc(len(items))
	if options.preAlloc > len(items) {
		preAlloc = HeapPreAlloc(options.preAlloc)
	}

	h := NewHeap[T](order, preAlloc)
	for i := range items {
		h.Push(items[i])
	}
	return h
}

func (h *Heap[T]) Push(item T) {
	h.Items = append(h.Items, item)
	h.heapifyUp(h.Len() - 1)
}

// Pop removes and returns the root item. The bool return is false if the
// heap was empty, in which case the returned T is the zero value.
func (h *Heap[T]) Pop() (T, bool) {
	if h.Len() == 0 {
		var zero T
		return zero, false
	}

	root := h.Items[0]
	lastIdx := h.Len() - 1

	h.Items[0] = h.Items[lastIdx]
	h.Items = h.Items[:lastIdx]

	h.heapifyDown(0)

	return root, true
}

func (h *Heap[T]) Len() int {
	return len(h.Items)
}

func (h *Heap[T]) Clone() Heap[T] {
	cloned := make([]T, h.Len())
	copy(cloned, h.Items)
	return Heap[T]{
		Items:    cloned,
		LessFunc: h.LessFunc,
	}
}

// Slice returns the items as sorted slice.
// Can be called many times with 0 changes to h.
func (h *Heap[T]) Slice() []T {
	clone := h.Clone()
	slice := make([]T, h.Len())
	for i := range clone.Len() {
		slice[i] = clone.PopValue()
	}
	return slice
}

// Drain returns the items as sorted slice.
// The return value is the same as with [Heap.Slice], but Drain consumes the whole of h.
func (h *Heap[T]) Drain() []T {
	slice := make([]T, h.Len())
	for i := range h.Len() {
		slice[i] = h.PopValue()
	}
	return slice
}

func (h *Heap[T]) IsEmpty() bool {
	return len(h.Items) == 0
}

// Peek returns the root item without removing it. The bool return is false
// if the heap is empty, in which case the returned T is the zero value.
func (h *Heap[T]) Peek() (T, bool) {
	if len(h.Items) == 0 {
		var zero T
		return zero, false
	}
	return h.Items[0], true
}

// PopValue is like Pop, but returns the zero value of T instead of a bool
// when the heap is empty.
func (h *Heap[T]) PopValue() T {
	value, _ := h.Pop()
	return value
}

// PeekValue is like Peek, but returns the zero value of T instead of a bool
// when the heap is empty.
func (h *Heap[T]) PeekValue() T {
	value, _ := h.Peek()
	return value
}

func (h *Heap[T]) heapifyUp(from int) {
	curr := from
	for curr != 0 {
		parent := ParentIdx(curr)
		if !h.LessFunc(h.Items, curr, parent) {
			break
		}

		h.swap(curr, parent)
		curr = parent
	}
}

func (h *Heap[T]) heapifyDown(from int) {
	curr := from
	length := len(h.Items)

	for {
		childLeft := LeftChildIdx(curr)
		if childLeft >= length {
			break
		}

		// Choose left child if:
		// 1) right child is null (out of range)
		// 2) left child has higher priority (lessFunc -> true)
		//
		// Otherwise use right child
		childRight := RightChildIdx(curr)
		child := childRight
		switch {
		case
			childRight >= length,
			h.LessFunc(h.Items, childLeft, childRight):
			child = childLeft
		}
		if h.LessFunc(h.Items, curr, child) {
			break
		}
		h.swap(curr, child)
		curr = child
	}
}

func (h *Heap[T]) swap(i, j int) {
	h.Items[i], h.Items[j] = h.Items[j], h.Items[i]
}

func parseOptions(opts ...HeapOption) HeapOptions {
	var options HeapOptions
	for i := range opts {
		opts[i](&options)
	}
	return options
}
