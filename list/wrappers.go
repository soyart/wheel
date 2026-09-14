package list

import "github.com/soyart/wheel"

// WrappedList is a BasicList that is wrapped inside another BasicList implementation.
type WrappedList[T any, L BasicList[T]] BasicList[T]

// SafeList is used as parameter in function where concurrency is used.
type SafeList[T any, L BasicList[T]] WrappedList[T, L]

// SetList is used as parameter in function where you'll need characteristics of a set.
type SetList[T comparable, L BasicList[T]] interface {
	WrappedList[T, L]
	wheel.Set[T]
}

// SafeSetList is a [SafeList] with a backing of a [SetList].
// It provides concurrency safety guarantee on top of set functionality (ie no duplicates)
//
// TODO: WTF?
type SafeSetList[T comparable, L BasicList[T]] SafeList[T, SetList[T, L]]
