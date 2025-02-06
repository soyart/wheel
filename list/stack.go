package list

type StackImpl[T any] []T

func NewStack[T any]() *StackImpl[T] {
	return new(StackImpl[T])
}

func NewStackSafe[T any]() SafeList[T, *StackImpl[T]] {
	return WrapSafeList[T](new(StackImpl[T]))
}

func (s *StackImpl[T]) Push(x T) {
	*s = append(*s, x)
}

func (s *StackImpl[T]) PushSlice(slice []T) {
	for _, elem := range slice {
		s.Push(elem)
	}
}

// Pop pops and returns the right-most element of self,
// returning nil if self is empty
func (s *StackImpl[T]) Pop() *T {
	state := *s
	l := len(state)
	if l == 0 {
		return nil
	}

	var elem T
	elem, *s = state[l-1], state[0:l-1]

	return &elem
}

func (s *StackImpl[T]) Len() int {
	state := *s
	return len(state)
}

func (s *StackImpl[T]) IsEmpty() bool {
	state := *s
	return len(state) == 0
}
