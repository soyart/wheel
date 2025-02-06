package wheel

import (
	"math/big"
	"testing"
)

type testSort[T any] struct {
	inputs   []T
	expected []T
	ordering SortOrder
}

type testSortCmp[T CmpOrdered[T]] testSort[T]

type foo struct {
	value int64
}

func (f *foo) Cmp(other *foo) int {
	switch {
	case f.value == other.value:
		return 0
	case f.value > other.value:
		return 1
	}

	return -1
}

func testCasesDefault() []testSort[int64] {
	return []testSort[int64]{
		{
			inputs:   []int64{},
			expected: []int64{},
			ordering: Ascending,
		},
		{
			inputs:   []int64{-1},
			expected: []int64{-1},
			ordering: Ascending,
		},
		{
			inputs:   []int64{2, 2, 2},
			expected: []int64{2, 2, 2},
			ordering: Ascending,
		},
		{
			inputs:   []int64{2, 3, 60, 1, 70, 234, -1},
			expected: []int64{-1, 1, 2, 3, 60, 70, 234},
			ordering: Ascending,
		},
		{
			inputs:   []int64{2, 3, 60, 1, 70, 234, -1},
			expected: []int64{234, 70, 60, 3, 2, 1, -1},
			ordering: Descending,
		},
	}
}

func testCasesBigInt() []testSortCmp[*big.Int] {
	testCasesInt64 := testCasesDefault()

	tests := make([]testSortCmp[*big.Int], len(testCasesInt64))
	for i := range testCasesInt64 {
		testCaseInt64 := &testCasesInt64[i]
		inputs := make([]*big.Int, len(testCaseInt64.inputs))
		expected := make([]*big.Int, len(testCaseInt64.expected))

		for j := range inputs {
			inputs[j] = big.NewInt(testCaseInt64.inputs[j])
		}
		for j := range expected {
			expected[j] = big.NewInt(testCaseInt64.expected[j])
		}

		tests[i] = testSortCmp[*big.Int]{
			inputs:   inputs,
			expected: expected,
			ordering: testCaseInt64.ordering,
		}
	}

	return tests
}

func testCasesFoo() []testSortCmp[*foo] {
	testCasesInt64 := testCasesDefault()
	tests := make([]testSortCmp[*foo], len(testCasesInt64))

	for i := range testCasesInt64 {
		testCaseInt64 := testCasesInt64[i]
		inputs := make([]*foo, len(testCaseInt64.inputs))
		expected := make([]*foo, len(testCaseInt64.expected))

		for j := range testCaseInt64.inputs {
			inputs[j] = &foo{
				value: testCaseInt64.inputs[j],
			}
		}
		for j := range testCaseInt64.expected {
			expected[j] = &foo{
				value: testCaseInt64.expected[j],
			}
		}

		tests[i] = testSortCmp[*foo]{
			inputs:   inputs,
			expected: expected,
			ordering: testCaseInt64.ordering,
		}
	}

	return tests
}

func testSortFunc[T comparable](t *testing.T, testCase testSort[T], sortFunc func([]T, SortOrder) []T) {
	result := sortFunc(testCase.inputs, testCase.ordering)

	for i := range testCase.expected {
		expected := testCase.expected[i]
		actual := result[i]

		if expected != actual {
			t.Fatalf("got unexpected value at position %d: expecting %+v, got %+v", i, expected, actual)
		}
	}
}

func testSortFuncCmp[T CmpOrdered[T]](t *testing.T, testCase testSortCmp[T], sortFunc func([]T, SortOrder) []T) {
	result := sortFunc(testCase.inputs, testCase.ordering)

	for i := range testCase.expected {
		expected := testCase.expected[i]
		actual := result[i]

		if cmp := expected.Cmp(actual); cmp != 0 {
			t.Logf("unexpected value at position %d - expecting %v, got %v", i, expected, actual)
			t.Fatalf("unexpected Cmp result - expecting 0, got %d", cmp)
		}
	}
}
