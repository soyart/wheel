package wheel

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	tests := testCasesDefault()
	for i := range tests {
		testCase := tests[i]
		testSortFunc(t, testCase, MergeSort)
	}

	// See if it'll overflow on 10M ints
	var s []int = make([]int, 10000000)
	for i := 0; i < 10000000; i++ {
		s[i] = i
	}

	MergeSort(s, Descending)
}

func TestMergeSortBigInts(t *testing.T) {
	tests := testCasesBigInt()
	for i := range tests {
		testSortFuncCmp(t, tests[i], MergeSortCmp)
	}
}

func TestMergeSortCustomCmp(t *testing.T) {
	tests := testCasesFoo()
	for i := range tests {
		testSortFuncCmp[*foo](t, tests[i], MergeSortCmp)
	}
}
