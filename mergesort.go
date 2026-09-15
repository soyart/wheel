package wheel

import (
	"cmp"
)

func MergeSort[T cmp.Ordered](arr []T, ordering SortOrder) []T {
	length := len(arr)

	if length < 2 {
		return arr
	}

	if length == 2 {
		a, b := arr[0], arr[1]

		lessFunc := LessFuncOrdered[T](ordering)

		if lessFunc(a, b) {
			return arr
		}

		return []T{b, a}
	}

	mid := length / 2

	left := MergeSort(arr[:mid], ordering)
	right := MergeSort(arr[mid:], ordering)

	return MergeSortedArrays(left, right, LessFuncOrdered[T](ordering))
}

func MergeSortCmp[T CmpOrdered[T]](arr []T, ordering SortOrder) []T {
	length := len(arr)

	if length < 2 {
		return arr
	}

	cmpLess := -1
	if ordering == Descending {
		cmpLess = 1
	}

	if length == 2 {
		a, b := arr[0], arr[1]

		if a.Cmp(b) == cmpLess {
			return arr
		}

		return []T{b, a}
	}

	mid := length / 2

	left := MergeSortCmp(arr[:mid], ordering)
	right := MergeSortCmp(arr[mid:], ordering)

	return MergeSortedArrays(left, right, func(a, b T) bool {
		return a.Cmp(b) == cmpLess
	})
}

func MergeSortedArrays[T any](
	a []T,
	b []T,
	lessFunc func(a T, b T) bool,
) []T {
	var p, ap, bp int
	sorted := make([]T, len(a)+len(b))

	for ap < len(a) && bp < len(b) {
		elemA, elemB := a[ap], b[bp]

		if lessFunc(elemA, elemB) {
			sorted[p] = elemA
			ap++
			p++

			continue
		}

		sorted[p] = elemB
		bp++
		p++
	}

	for ap < len(a) {
		sorted[p] = a[ap]
		ap++
		p++
	}

	for bp < len(b) {
		sorted[p] = b[bp]
		bp++
		p++
	}

	return sorted
}
