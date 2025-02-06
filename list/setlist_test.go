package list

import "testing"

func TestSet(t *testing.T) {
	data := []uint16{1, 2, 1, 2, 3, 4, 2, 1}
	set := ToSetList(data)

	counts := make(map[uint16]int)
	for !set.IsEmpty() {
		popped := set.Pop()
		if popped == nil {
			t.Fatal("popped nil - should not happen")
		}
		val := *popped
		counts[val]++
	}

	for _, item := range data {
		if counts[item] != 1 {
			t.Fatal("set failed")
		}
	}

	var item uint16 = 69
	for i := 0; i < 10; i++ {
		set.Push(item)
	}
	if set.Len() != 1 {
		t.Fatal("set failed")
	}
}
