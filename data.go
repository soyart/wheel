package wheel

import (
	"cmp"
)

type Getter[T any] interface {
	GetValue() T
}

type Set[T comparable] interface {
	HasDuplicate(x T) bool
}

type wrapper[T any] struct {
	inner T
}

func (w *wrapper[T]) SetValue(value T) {
	w.inner = value
}

func (w *wrapper[T]) GetValue() T {
	return w.inner
}

func NewGetter[T any](t T) Getter[T] {
	return &wrapper[T]{inner: t}
}

func ToValues[T any](getters []Getter[T]) []T {
	values := make([]T, len(getters))
	for i := range getters {
		values[i] = getters[i].GetValue()
	}

	return values
}

func FromValues[T any](values []T) []Getter[T] {
	getters := make([]Getter[T], len(values))
	for i := range getters {
		values[i] = getters[i].GetValue()
	}

	return getters
}

func MaxValuer[T cmp.Ordered](values []Getter[T]) T {
	var t T
	if len(values) == 0 {
		return t
	}

	max := values[0].GetValue()
	for i := range values {
		v := values[i].GetValue()
		if v > max {
			max = v
		}
	}

	return max
}

func MinValuer[T cmp.Ordered](values []Getter[T]) T {
	var t T
	if len(values) == 0 {
		return t
	}

	min := values[0].GetValue()
	for i := range values {
		v := values[i].GetValue()
		if v < min {
			min = v
		}
	}

	return min
}
