package tree

import (
	"cmp"

	"github.com/soyart/wheel"
)

// Heap is a binary heap implementation backed by Go slice.
type Heap[T any] struct {
	Items    []wheel.Getter[T]
	LessFunc wheel.LessFunc[wheel.Getter[T]]
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
	options := parseOptions(opts...)
	return &Heap[T]{
		Items:    make([]wheel.Getter[T], 0, options.preAlloc),
		LessFunc: wheel.FactoryLessFuncOrdered[T](order),
	}
}

func NewHeapCmp[T wheel.CmpOrdered[T]](order wheel.SortOrder, opts ...HeapOption) *Heap[T] {
	options := parseOptions(opts...)
	return &Heap[T]{
		Items:    make([]wheel.Getter[T], 0, options.preAlloc),
		LessFunc: wheel.FactoryLessFuncCmp[T](order),
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
	getter := wheel.NewGetter(item)
	h.PushGetter(getter)
}

func (h *Heap[T]) PushGetter(getter wheel.Getter[T]) {
	h.Items = append(h.Items, getter)
	h.heapifyUp(h.Len() - 1)
}

func (h *Heap[T]) Pop() *T {
	root := h.PopGetter()
	if root == nil {
		return nil
	}

	rootValue := root.GetValue()
	return &rootValue
}

func (h *Heap[T]) PopGetter() wheel.Getter[T] {
	if h.Len() == 0 {
		return nil
	}

	rootNode := h.Items[0]
	lastIdx := h.Len() - 1

	h.Items[0] = h.Items[lastIdx]
	h.Items = h.Items[:lastIdx]

	h.heapifyDown(0)

	return rootNode
}

func (h *Heap[T]) Len() int {
	return len(h.Items)
}

func (h *Heap[T]) Clone() Heap[T] {
	cloned := make([]wheel.Getter[T], h.Len())
	for i := range cloned {
		cloned[i] = h.Items[i]
	}
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

func (h *Heap[T]) PeekGetter() wheel.Getter[T] {
	if len(h.Items) == 0 {
		return nil
	}
	return h.Items[0]
}

func (h *Heap[T]) PopValue() T {
	copied := *h.Pop()
	return copied
}

func (h *Heap[T]) PeekValue() T {
	if getter := h.PeekGetter(); getter != nil {
		return getter.GetValue()
	}
	var zero T
	return zero
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
