package list

import "github.com/soyart/wheel"

type WrappedList[T any, L BasicList[T]] BasicList[T]

// Use SafeList as parameter in function where concurrency is used.
type SafeList[T any, L BasicList[T]] WrappedList[T, L]

// Use SetList as parameter in function where you'll need characteristics of a set.
type SetList[T comparable, L BasicList[T]] interface {
	WrappedList[T, L]
	wheel.Set[T]
}

// TODO: WTF?
type SafeSetList[T comparable, L BasicList[T]] SafeList[T, SetList[T, L]]
