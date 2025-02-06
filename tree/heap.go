package tree

import (
	"golang.org/x/exp/constraints"

	"github.com/soyart/wheel"
)

// Heap is a binary heap implementation backed by Go slice.
type Heap[T any] struct {
	Items    []wheel.Getter[T]
	LessFunc wheel.LessFunc[wheel.Getter[T]]
}

func NewHeap[T constraints.Ordered](
	order wheel.SortOrder,
) *Heap[T] {
	return &Heap[T]{
		LessFunc: wheel.FactoryLessFuncOrdered[T](order),
	}
}

func NewHeapCmp[T wheel.CmpOrdered[T]](
	order wheel.SortOrder,
) *Heap[T] {
	return &Heap[T]{
		LessFunc: wheel.FactoryLessFuncCmp[T](order),
	}
}

func (h *Heap[T]) Push(item T) {
	getter := wheel.NewGetter[T](item)
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

	var t T
	return t
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

		childRight := RightChildIdx(curr)

		// Child to compare
		//nolint:ineffassign
		child := -1

		// Choose left child if:
		// 1) right child is null (out of range)
		// 2) left child has higher priority (lessFunc -> true)
		//
		// Otherwise use right child
		switch {
		case
			childRight >= length,
			h.LessFunc(h.Items, childLeft, childRight):

			child = childLeft

		default:
			child = childRight
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
