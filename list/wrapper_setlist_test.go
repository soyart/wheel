package list

import "testing"

// This function will take in any BasicList[T]
func takeAnyList[T any](l BasicList[T]) {}

// If you want your parameter to be a safe stack
func takeSafeStack[T comparable](ss SafeList[T, *StackImpl[T]]) {}

// If you want your parameter to be a safe queue
func takeSafeQueue[T comparable](sl SafeList[T, *QueueImpl[T]]) {}

// If you want your parameter to be any set list
func takeSetList[T comparable](sl SetList[T, BasicList[T]]) {}

// If you want your parameter to be SafeList[T, *SetListImpl[T]] (a *SetListImpl[T] wrapped in a SafeList[T, L])
// TODO: 1 problem with SafeList[T, *SetListImpl[T]] - it cannot call HasDuplicate!
// func takeSafeSetList[T comparable](SafeList[T, *SetListImpl[T]]) {}
func takeSafeSetList[T comparable](ssl SafeSetList[T, *SetListImpl[T]]) {
	// var zeroValue T
	// ssl.HasDuplicate(zeroValue)
}

func testCompileSetList(t *testing.T) {
	setList := ToSetList([]int{1, 2, 1, 2})
	takeSetList[int](setList)
	safeSetList := WrapSafeList[int](setList)
	takeSafeSetList[int](safeSetList)

	testSetListQueue(t, []float64{1, 2, 4, 2, 4, 1, 3, 5})
	testSetListQueue(t, []string{"foo", "bar", "baz", "bar", "foom"})
	testSetListStack(t, []float64{1, 2, 4, 2, 4, 1, 3, 5})
	testSetListStack(t, []string{"foo", "bar", "baz", "bar", "foom"})
}

func testSetListStack[T comparable](t *testing.T, data []T) {
	// Standard SetListImpl
	set := ToSetList(data)

	// Wrap stack with SetListWrapper and then wrap that shit with SafeListWrapper
	stack := NewStack[T]()                     // *StackImpl[T]
	setStack := WrapSetListKeepInner[T](stack) // *SetList[T, *Stack[T]]
	safeSetStack := WrapSafeList[T](setStack)  // *SafeList[T, *SetListWrapper[T]]

	// Wrap anotherStack with SafeListWrapper and then wrap that shit with SetListWrapper
	anotherStack := NewStack[T]()                      // *Stack[T]
	safeStack := WrapSafeList[T](anotherStack)         // *SafeList[T, *Stack[T]]
	setSafeStack := WrapSetListKeepInner[T](safeStack) // *SetList[T, *SafeList[T, Stack[T]]]

	// All 6 should implement BasicList[T] and Stack[T]
	lists := []BasicList[T]{
		set,
		stack, setStack, safeSetStack,
		anotherStack, safeStack, setSafeStack,
	}
	for _, l := range lists {
		takeAnyList(l)
	}

	// If wrapped lastly be SafeList[T, BasicList[T]] (like safeSetStack), then it's obviously NOT Set[T].
	sets := []SetList[T, BasicList[T]]{
		set,
		setStack,
		// safeSetStack, // Compile error!
		setSafeStack,
	}
	for _, x := range sets {
		takeSetList(x)
	}

	safes := []SafeList[T, *StackImpl[T]]{
		safeSetStack,
		safeStack,
		// setSafeStack, // Compile error!
	}
	for _, x := range safes {
		takeSafeStack(x)
	}

	testSets := []SetList[T, BasicList[T]]{set, setStack, setSafeStack}
	for _, set := range testSets {
		testSetPushAndPop[T](t, set, data)
	}
}

func testSetListQueue[T comparable](t *testing.T, data []T) {
	set := ToSetList(data)

	queue := NewQueue[T]()                     // *QueueImpl[T]
	setQueue := WrapSetListKeepInner[T](queue) // *SetListWrapper[T, *QueueImpl[T]]
	safeSetQueue := WrapSafeList[T](setQueue)  // *SafeListWrapper[T, *SetListWrapper[T, *QueueImpl]]

	anotherQueue := NewQueue[T]()                      // *QueueImpl[T]
	safeQueue := WrapSafeList[T](queue)                // *SafeListWrapper[T, *QueueImpl[T]]
	setSafeQueue := WrapSetListKeepInner[T](safeQueue) // *SetListWrapper[T, SafeListWrapper[T, *QueueImpl[T]]]

	// All 6 should implement BasicList[T] and Queue[T]
	lists := []BasicList[T]{
		set,
		queue, setQueue, safeSetQueue,
		anotherQueue, safeQueue, setSafeQueue,
	}
	for _, l := range lists {
		takeAnyList(l)
	}

	// If wrapped lastly be SafeList[T, BasicList[T]] (like safeSetQueue), then it's obviously NOT Set[T].
	sets := []SetList[T, BasicList[T]]{
		set,
		setQueue,
		// safeSetQueue, // Compile error!
		setSafeQueue,
	}
	for _, x := range sets {
		takeSetList(x)
	}

	safes := []SafeList[T, *QueueImpl[T]]{
		safeQueue,
		safeSetQueue,
		// setSafeQueue, // Compile error!
	}
	for _, x := range safes {
		takeSafeQueue(x)
	}

	testSets := []SetList[T, BasicList[T]]{set, setQueue, setSafeQueue}
	for _, set := range testSets {
		testSetPushAndPop[T](t, set, data)
	}
}

func testSetPushAndPop[T comparable](t *testing.T, setList BasicList[T], data []T) {
	for _, item := range data {
		setList.Push(item)
	}

	seenCounts := make(map[T]int)
	for !setList.IsEmpty() {
		popped := setList.Pop()
		if popped == nil {
			t.Fatal("popped nil - should not happen")
		}

		value := *popped
		seenCounts[value]++
	}

	for value, count := range seenCounts {
		if count != 1 {
			t.Fatalf("has duplicates for key %v in set\n", value)
		}
	}
}
