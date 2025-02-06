package wheel

import "golang.org/x/exp/constraints"

// CopySlice return a copy of |arr|.
func CopySlice[T any](arr []T) []T {
	ret := make([]T, len(arr))
	copy(ret, arr)
	return ret
}

// ReverseInPlace reverses |arr| in-place.
func ReverseInPlace[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// Reverse returns a new slice, with elements being those of |arr| after being reversed.
func Reverse[T any](arr []T) []T {
	reversed := CopySlice(arr)
	ReverseInPlace(reversed)
	return reversed
}

// Contains iterates over |arr| and returns true if an element in |arr| is equal to |item|.
func Contains[T comparable](arr []T, item T) bool {
	for _, elem := range arr {
		if elem == item {
			return true
		}
	}
	return false
}

// FilterSlice iterates over |arr| and calls |filterFunc| on each of the element |elem|.
// If filterFunc returns true, |elem| is added to the return slice.
func FilterSlice[T any](arr []T, filterFunc func(elem T) bool) []T {
	if arr == nil {
		return nil
	}

	var filtered []T
	for _, elem := range arr {
		if filterFunc(elem) {
			filtered = append(filtered, elem)
		}
	}
	return filtered
}

// CollectPointers iterates over |arr| and returns a new slice containing all references
// to the elements of |arr|.
func CollectPointers[T any](arr []T) []*T {
	if arr == nil {
		return nil
	}

	out := make([]*T, len(arr))
	for i := range arr {
		out[i] = &arr[i]
	}
	return out
}

// CollectPointersIf iterates over |arr| and calling filterFunc on the element |elem|.
// If filterFunc returns true, |elem| is added to the return slice.
func CollectPointersIf[T any](arr []T, filterFunc func(elem T) bool) []*T {
	if arr == nil {
		return nil
	}

	var filtered []*T
	for i := range arr {
		value := arr[i]
		if filterFunc(value) {
			filtered = append(filtered, &value)
		}
	}
	return filtered
}

// DerefValues iterates over |arr| and return a new slice containing dereferenced values of the elements in |arr|.
// If the pointer at `arr[i]` is nil, the returned slice will have default value of type `T` at index `i`.
func DerefValues[T any](arr []*T) []T {
	if arr == nil {
		return nil
	}

	var t T
	values := make([]T, len(arr))
	for i := range arr {
		if arr[i] == nil {
			values[i] = t
			continue
		}
		values[i] = *arr[i]
	}

	return values
}

// DerefValuesIf iterates over |arr| and calls |filterFunc| on the element |elem|.
// If filterFunc returns true, |elem|'s dereferenced value is added to the return slice.
// Unlike DerefValues but similar to CollectPointersIf, if |elem| is nil,
// its value will not be added to the return slice
func DerefValuesIf[T any](arr []*T, filterFunc func(elem T) bool) []T {
	if arr == nil {
		return nil
	}

	var filtered []T
	for _, elem := range arr {
		if elem == nil {
			continue
		}
		value := *elem
		if filterFunc(value) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

// Map maps |arr| to a new slice using |mapFunc|. |mapFunc| returns a tuple of generic type U and a boolean.
// If |mapFunc|'s 2nd return value is true, its 1st return value (of type U) will be appended to the result slice.
// Note: `T` and `U` can be of the same type.
func Map[T any, U any](
	arr []T,
	mapFunc func(elem T) (U, bool),
) []U {
	if arr == nil {
		return nil
	}

	var mapped []U //nolint:prealloc
	for _, elem := range arr {
		u, ok := mapFunc(elem)
		if !ok {
			continue
		}
		mapped = append(mapped, u)
	}

	return mapped
}

func Max[T constraints.Ordered](items ...T) T {
	if len(items) == 0 {
		var zeroValue T
		return zeroValue
	}

	max := items[0]
	for _, num := range items {
		if num > max {
			max = num
		}
	}
	return max
}

func Min[T constraints.Ordered](items ...T) T {
	if len(items) == 0 {
		var zeroValue T
		return zeroValue
	}

	min := items[0]
	for _, num := range items {
		if num < min {
			min = num
		}
	}
	return min
}

func Sum[T constraints.Integer | constraints.Float](items ...T) T {
	var sum T
	for _, item := range items {
		sum += item
	}
	return sum
}

func Avg[T constraints.Integer | constraints.Float](items ...T) T {
	return Sum(items...) / T(len(items))
}
