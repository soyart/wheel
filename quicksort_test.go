package wheel

import (
	"testing"
)

func TestQuickSortNoCopy(t *testing.T) {
	tests := [][]int{
		{3, 1, 9, 2, 1, 8, 4},
		{10, 40, 30, 80, 20},
		{100, 200, 300, 400},
	}

	for i := range tests {
		data := tests[i]
		QuickSortNoCopy(data, Ascending, 0, len(data)-1)

		prev := data[0]
		for j := range data {
			elem := data[j]
			if elem < prev {
				t.Fatal("unexpected elem < prev")
			}
		}
	}
}

func TestQuickSort(t *testing.T) {
	tests := testCasesDefault()
	for i := range tests {
		testCase := tests[i]
		testSortFunc(t, testCase, QuickSort)
	}

	// See if it'll overflow on 10M ints
	var s []int = make([]int, 10000000)
	for i := 0; i < 10000000; i++ {
		s[i] = i
	}

	QuickSort(s, Descending)
}

func TestQuickSortBigInts(t *testing.T) {
	tests := testCasesBigInt()
	for i := range tests {
		testSortFuncCmp(t, tests[i], QuickSortCmp)
	}
}

func TestQuickSortCustomCmp(t *testing.T) {
	tests := testCasesFoo()
	for i := range tests {
		testSortFuncCmp[*foo](t, tests[i], QuickSortCmp)
	}
}
