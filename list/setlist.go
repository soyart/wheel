package list

// SetListImpl implements SetList[T, L] and BasicList[T]
//
// A *set list* is a list (implements BasicList[T]) that has characteristics of a set.
// No duplicate is allowed in a set list.
//
// If you want set functionality for other types, e.g., for QueueImpl[T],
// then you can wrap the queue using WrapSetList(l BasicList[T]) to get a SetList[T, *QueueImpl[T]].
type SetListImpl[T comparable] struct {
	haystack   []T
	duplicates map[T]struct{}
	length     int
}

func NewSetList[T comparable]() *SetListImpl[T] {
	return new(SetListImpl[T])
}

// ToSetList iterates over src T and push it to a new SetListImpl[T]
func ToSetList[T comparable](src []T) SetList[T, *SetListImpl[T]] {
	var haystack []T
	duplicates := make(map[T]struct{})
	var length int
	for _, item := range src {
		if _, found := duplicates[item]; !found {
			haystack = append(haystack, item)
			duplicates[item] = struct{}{}
			length++
		}
	}

	return &SetListImpl[T]{
		haystack:   haystack,
		duplicates: duplicates,
		length:     length,
	}
}

func (s *SetListImpl[T]) HasDuplicate(x T) bool {
	_, found := s.duplicates[x]
	return found
}

func (s *SetListImpl[T]) Push(x T) {
	if !s.HasDuplicate(x) {
		s.haystack = append(s.haystack, x)
		s.duplicates[x] = struct{}{}
		s.length++
	}
}

func (s *SetListImpl[T]) PushSlice(slice []T) {
	for _, elem := range slice {
		s.Push(elem)
	}
}

func (s *SetListImpl[T]) Pop() *T {
	toPop := s.haystack[s.length-1]
	s.length--
	return &toPop
}

func (s *SetListImpl[T]) Len() int {
	return s.length
}

func (s *SetListImpl[T]) IsEmpty() bool {
	return s.length == 0
}
