package list

import "sync"

// SafeListWrapper[T] wraps BasicList[T] and uses sync.RWMutex to avoid data races.
// L was added to make sure that the underlying list type is always accessible from the instance type,
// for example, SafeListWrapper[float64, *Queue[float64]]
// if L was not the type parameter, then a safe uint8 stack, safe uint8 queue, etc,
// will be indifferentiable with the same type `SafeListWrapper[uint8]`.
// SafeListWrapper[T, BasicList[T]] also implements BasicList[T]
type SafeListWrapper[T any, L BasicList[T]] struct {
	inner L
	mut   *sync.RWMutex
}

// WrapSafeList[T] wraps a BasicList[T] into SafeListWrapper[T],
// where T is the underlying entity (item) type and L is the underlying BasicList[T] type.
// If you're wrapping a variable `foo“ of type `*Stack[uint8]`, then call this function with:
// WrapSafeList[uint8, *Stack[uint8]](foo)
func WrapSafeList[T any, L BasicList[T]](basicList L) *SafeListWrapper[T, L] {
	return &SafeListWrapper[T, L]{
		inner: basicList,
		mut:   &sync.RWMutex{},
	}
}

func (w *SafeListWrapper[T, L]) Push(x T) {
	w.mut.Lock()
	defer w.mut.Unlock()

	w.inner.Push(x)
}

func (w *SafeListWrapper[T, L]) PushSlice(x []T) {
	w.mut.Lock()
	defer w.mut.Unlock()

	w.inner.PushSlice(x)
}

func (w *SafeListWrapper[T, L]) Pop() *T {
	w.mut.Lock()
	defer w.mut.Unlock()

	return w.inner.Pop()
}

func (w *SafeListWrapper[T, L]) Len() int {
	w.mut.RLock()
	defer w.mut.RUnlock()

	return w.inner.Len()
}

func (w *SafeListWrapper[T, L]) IsEmpty() bool {
	w.mut.RLock()
	defer w.mut.RUnlock()

	return w.inner.IsEmpty()
}

func (w *SafeListWrapper[T, L]) IsSafe() bool {
	return true
}
