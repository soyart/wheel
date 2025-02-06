package list

import "testing"

func testSafeListWrapper(t *testing.T) {
	valuesComposite := []interface{}{make(chan int), []byte("test"), &struct{ n int }{n: 69}}
	ints := []int{1, 2, 3, 4}
	floats := []float32{1, 2, 3, 4}
	strings := []string{"kuy", "hee", "tad"}

	testSafeList(t, valuesComposite)
	testSafeList(t, ints)
	testSafeList(t, floats)
	testSafeList(t, strings)
}

func testSafeList[T any](t *testing.T, values []T) {
	// Those list types wrapped in SafeList[T, BasicList[T]] also implement BasicList[T]
	basicStack := NewStackSafe[T]()
	safeStack := WrapSafeList[T](basicStack)
	basicQueue := NewQueue[T]()
	safeQueue := WrapSafeList[T](basicQueue)
	anotherSafeQueue := NewQueueSafe[T]()
	tests := []BasicList[T]{safeStack, safeQueue, anotherSafeQueue}

	for _, l := range tests {
		testBasicList(t, values, l)
	}
}
