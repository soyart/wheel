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
) LessFunc[Getter[T]] {
	if order == Ascending {
		return func(items []Getter[T], i, j int) bool {
			return items[i].GetValue() < items[j].GetValue()
		}
	}

	return func(items []Getter[T], i, j int) bool {
		return items[i].GetValue() > items[j].GetValue()
	}
}

func FactoryLessFuncCmp[T CmpOrdered[T]](
	order SortOrder,
) LessFunc[Getter[T]] {
	if order == Ascending {
		return func(items []Getter[T], i, j int) bool {
			return items[i].GetValue().Cmp(items[j].GetValue()) < 0
		}
	}

	return func(items []Getter[T], i, j int) bool {
		return items[i].GetValue().Cmp(items[j].GetValue()) > 0
	}
}
