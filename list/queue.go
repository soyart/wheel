package list

type QueueImpl[T any] []T

func NewQueue[T any]() *QueueImpl[T] {
	return new(QueueImpl[T])
}

func NewQueueSafe[T any]() SafeList[T, *QueueImpl[T]] {
	return WrapSafeList[T](new(QueueImpl[T]))
}

func (s *QueueImpl[T]) Push(x T) {
	*s = append(*s, x)
}

func (s *QueueImpl[T]) PushSlice(slice []T) {
	for _, elem := range slice {
		s.Push(elem)
	}
}

// Pop pops and returns the left-most element of self,
// returning nil if self is empty
func (s *QueueImpl[T]) Pop() *T {
	state := *s
	l := len(state)
	if l == 0 {
		return nil
	}

	var elem T
	elem, *s = state[0], state[1:l]

	return &elem
}

func (s *QueueImpl[T]) Len() int {
	state := *s
	return len(state)
}

func (s *QueueImpl[T]) IsEmpty() bool {
	state := *s
	return len(state) == 0
}
