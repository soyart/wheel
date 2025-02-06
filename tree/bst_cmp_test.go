package tree

import (
	"math/big"
	"testing"
)

type wrapper int

func (w *wrapper) Cmp(other *wrapper) int {
	switch {
	case *w > *other:
		return 1

	case *w < *other:
		return -1
	}

	return 0
}

func TestBstCustomInsert(t *testing.T) {
	bst := NewBstCmp[*wrapper]()

	start := 1
	limit := 10

	for i := start; i < limit; i++ {
		w := wrapper(i)
		inserted := bst.Insert(&w)
		if !inserted {
			t.Fatalf("Insert returned false for new value %d", i)
		}
	}

	for i := start; i < limit; i++ {
		w := wrapper(i)
		inserted := bst.Insert(&w)
		if inserted {
			t.Fatalf("Insert returned false for new value %d", i)
		}
	}

	for i := start; i < limit; i++ {
		w := wrapper(i)
		if !bst.Find(&w) {
			t.Fatalf("missing node %d", i)
		}
	}

	outOfRange := wrapper(-1)
	if bst.Find(&outOfRange) {
		t.Fatalf("unexpected false positive for %d", outOfRange)
	}
}

func TestBstCustomRemoveEmpty(t *testing.T) {
	bst := new(BstCmp[*wrapper])

	items := []int{3, 1, 2, 0, 5}

	for i := range items {
		w := wrapper(items[i])
		bst.Insert(&w)

		t.Log("root after insert", items[i], "root", bst.Root)
	}

	for i := range items {
		w := wrapper(items[i])
		bst.Remove(&w)

		t.Log("root after delete", items[i], "root", bst.Root)
	}

	t.Log("final root", bst.Root)
	if !bst.Root.IsLeaf() {
		t.Fatalf("not all children removed")
	}

	if bst.Root.ok {
		t.Log("final root", bst.Root)
		t.Fatalf("root is still ok")
	}
}

func TestBstCustomRemove(t *testing.T) {
	bst := new(BstCmp[*wrapper])

	start := 1
	limit := 10
	target := 5
	targetWrapper := wrapper(target)

	for i := start; i < limit; i++ {
		w := wrapper(i)
		bst.Insert(&w)
	}

	if !bst.Remove(&targetWrapper) {
		t.Fatalf("remove returned false on target %d", target)
	}

	if bst.Find(&targetWrapper) {
		t.Fatalf("found removed target %d", target)
	}

	for i := start; i < limit; i++ {
		w := wrapper(i)
		if !bst.Find(&w) {
			if i == target {
				continue
			}

			t.Fatalf("missing node %d", i)
		}
	}
}

func TestBstCustomRemoveBigInt(t *testing.T) {
	bst := new(BstCmp[*big.Int])

	start := int64(1)
	limit := int64(10)
	target := int64(5)

	for i := start; i < limit; i++ {
		bst.Insert(big.NewInt(i))
	}

	if !bst.Remove(big.NewInt(target)) {
		t.Fatalf("remove returned false on target %d", target)
	}

	if bst.Find(big.NewInt(target)) {
		t.Fatalf("found removed target %d", target)
	}

	for i := start; i < limit; i++ {
		if !bst.Find(big.NewInt(i)) {
			if i == target {
				continue
			}

			t.Fatalf("missing node %d", i)
		}
	}
}
