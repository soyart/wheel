package list

// SetListWrapper wraps BasicList[T] into field `basicList`,
// and maintains a hash map of seen items and the basicList length
// so that determining duplicates and getting the length take O(1) time.
type SetListWrapper[T comparable, L BasicList[T]] struct {
	inner      L
	duplicates map[T]struct{}
}

// O(1)
func (s *SetListWrapper[T, L]) HasDuplicate(x T) bool {
	_, found := s.duplicates[x]
	return found
}

func (s *SetListWrapper[T, L]) Push(x T) {
	if s.HasDuplicate(x) {
		return
	}

	s.inner.Push(x)
	s.duplicates[x] = struct{}{}
}

func (s *SetListWrapper[T, L]) PushSlice(values []T) {
	for i := range values {
		s.Push(values[i])
	}
}

func (s *SetListWrapper[T, L]) Pop() *T {
	return s.inner.Pop()
}

func (s *SetListWrapper[T, L]) Len() int {
	return s.inner.Len()
}

func (s *SetListWrapper[T, L]) IsEmpty() bool {
	return s.inner.Len() == 0
}

func (s *SetListWrapper[T, L]) Inner() L {
	return s.inner
}

// WrapSetListKeepInner returns a SetListWrapper over inner.
// It does not remove duplicates in inner - only elements pushed thereafter
// are checked against the hash table.
func WrapSetListKeepInner[T comparable](inner BasicList[T]) *SetListWrapper[T, BasicList[T]] {
	return &SetListWrapper[T, BasicList[T]]{
		inner:      inner,
		duplicates: make(map[T]struct{}),
	}
}
