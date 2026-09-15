package wheel

import (
	"cmp"
)

func swap[T any](v []T, i, j int) {
	v[i], v[j] = v[j], v[i]
}

func QuickSortNoCopy[T cmp.Ordered](arr []T, ordering SortOrder, left, right int) {
	if left >= right {
		return
	}

	// [arr[left], ..., arr[mid], ...., arr[right]]
	mid := (left + right) / 2

	// [arr[mid], ..., arr[left], ...., arr[right]]
	swap(arr, left, mid) // Move our pivot (chosen as the mid element) to the left

	last := left
	for i := left + 1; i <= right; i++ {
		if arr[i] < arr[left] {
			last++
			swap(arr, i, last)
		}
	}

	// [pivot, ...elems<=pivot... , last, .... elems>=pivot ....]

	swap(arr, left, last) // Move pivot back

	QuickSortNoCopy(arr, ordering, left, last-1)
	QuickSortNoCopy(arr, ordering, last+1, right)
}

func QuickSort[T cmp.Ordered](arr []T, ordering SortOrder) []T {
	l := len(arr)
	// Base case
	if l < 2 {
		return arr
	}

	// Pivot here is going to be middle elem
	mid := l / 2
	pivot := arr[mid]

	// pivoted is the list without the |pivot| element/member
	pivoted := append(arr[:mid], arr[mid+1:]...) //nolint:gocritic
	var left, right []T                          //nolint:prealloc

	{
		// isLess should have lifetime of no more than the for-looop
		// to help minimize stack size in huge recursive calls
		isLess := LessFuncOrdered[T](ordering)
		for _, elem := range pivoted {
			if isLess(elem, pivot) {
				left = append(left, elem)
				continue
			}

			right = append(right, elem)
		}
	}

	sorted := append(QuickSort(left, ordering), pivot)
	sorted = append(sorted, QuickSort(right, ordering)...)

	return sorted
}

func QuickSortCmp[T CmpOrdered[T]](arr []T, ordering SortOrder) []T {
	l := len(arr)
	// Base case
	if l < 2 {
		return arr
	}

	// Pivot here is going to be middle elem
	mid := l / 2
	pivot := arr[mid]

	// pivoted is the list without the pivot element/member
	pivoted := append(arr[:mid], arr[mid+1:]...) //nolint: gocritic
	var left, right []T                          //nolint:prealloc

	cmpLess := -1
	if ordering == Descending {
		cmpLess = 1
	}

	for _, elem := range pivoted {
		if elem.Cmp(pivot) == cmpLess {
			left = append(left, elem)
			continue
		}

		right = append(right, elem)
	}

	sorted := append(QuickSortCmp(left, ordering), pivot)
	sorted = append(sorted, QuickSortCmp(right, ordering)...)

	return sorted
}
