package wheel

import (
	"cmp"
	"fmt"
)

type SortOrder uint8

const (
	Ascending SortOrder = iota
	Descending
)

// CmpOrdered represents any type T with `Cmp(T) int` method.
// Examples of types that implement this interface include *big.Int and *big.Float.
type CmpOrdered[T any] interface {
	Cmp(T) int
}

func badOrder(ordering SortOrder) string {
	return fmt.Sprintf("bad SortOrder %d", ordering)
}

func (d SortOrder) IsValid() bool {
	switch d {
	case Ascending, Descending:
		return true
	}

	return false
}

type LessFunc[T any] func(list []T, i, j int) bool

// LessFuncOrdered selects the appropriate comparison function to check if the elements are ordered.
// If the returned function returns true, then the elements are sorted according to its |ordering|
func LessFuncOrdered[T cmp.Ordered](ordering SortOrder) func(T, T) bool {
	switch ordering {
	case Ascending:
		return func(v1, v2 T) bool {
			return v1 <= v2
		}

	case Descending:
		return func(v1, v2 T) bool {
			return v1 >= v2
		}
	}

	panic(badOrder(ordering))
}

// Less implementation for constraints.Ordered
func FactoryLessFuncOrdered[T cmp.Ordered](
	order SortOrder,
) LessFunc[T] {
	if order == Ascending {
		return func(items []T, i, j int) bool {
			return items[i] < items[j]
		}
	}

	return func(items []T, i, j int) bool {
		return items[i] > items[j]
	}
}

func FactoryLessFuncCmp[T CmpOrdered[T]](
	order SortOrder,
) LessFunc[T] {
	if order == Ascending {
		return func(items []T, i, j int) bool {
			return items[i].Cmp(items[j]) < 0
		}
	}

	return func(items []T, i, j int) bool {
		return items[i].Cmp(items[j]) > 0
	}
}

// LessFuncBy builds a LessFunc[T] that orders T by a comparable key extracted
// via keyFunc. It is the explicit, allocation-free replacement for boxing T
// behind an interface just to get a uniform GetValue()-based comparator:
// useful when T is a rich type (e.g. a graph node) rather than an ordered
// value itself.
func LessFuncBy[T any, K cmp.Ordered](
	order SortOrder,
	keyFunc func(T) K,
) LessFunc[T] {
	if order == Ascending {
		return func(items []T, i, j int) bool {
			return keyFunc(items[i]) < keyFunc(items[j])
		}
	}

	return func(items []T, i, j int) bool {
		return keyFunc(items[i]) > keyFunc(items[j])
	}
}
