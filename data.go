package wheel

import (
	"cmp"
)

type Set[T comparable] interface {
	HasDuplicate(x T) bool
}

func MaxValuer[T cmp.Ordered](values []T) T {
	var t T
	if len(values) == 0 {
		return t
	}

	max := values[0]
	for i := range values {
		if values[i] > max {
			max = values[i]
		}
	}

	return max
}

func MinValuer[T cmp.Ordered](values []T) T {
	var t T
	if len(values) == 0 {
		return t
	}

	min := values[0]
	for i := range values {
		if values[i] < min {
			min = values[i]
		}
	}

	return min
}
