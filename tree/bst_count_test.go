package tree

import (
	"math/big"
	"testing"
)

func TestBstCountInsertFind(t *testing.T) {
	bst := new(Bst[int])
	bstCount := NewBstCount[int](bst)

	start := 1
	limit := 10

	var c uint64
	for i := start; i < limit; i++ {
		inserted := bstCount.Insert(i)
		if !inserted {
			t.Fatalf("Insert returned false after inserting new value %d", i)
		}

		c++
		if bstCount.Count != c {
			t.Fatalf("Unexpected BstCount count - expecting %d, got %d", c, bstCount.Count)
		}

		t.Logf("Root after inserted %d: %+v", i, bst.Root)
	}

	for i := start; i < limit; i++ {
		inserted := bstCount.Insert(i)
		if inserted {
			t.Fatalf("Insert returned true after replacing value %d with new node", i)
		}

		if !bst.Find(i) {
			t.Fatalf("missing node %d", i)
		}
	}

	if bstCount.Count != c {
		t.Fatalf("Unexpected BstCount count after reinsert duplicates - expecting %d, got %d", c, bstCount.Count)
	}

	outOfRange := -1
	if bst.Find(outOfRange) {
		t.Fatalf("unexpected false positive for %d", outOfRange)
	}

	curr := &bst.Root
	for curr != nil {
		t.Logf("node %+v", curr)
		curr = curr.right
	}
}

func TestBstCountRemove(t *testing.T) {
	bst := new(Bst[int])
	bstCount := NewBstCount(bst)

	start := 1
	limit := 12
	target := 5

	var c uint64
	for i := start; i < limit; i++ {
		bstCount.Insert(i)
		c++
	}

	if !bstCount.Remove(target) {
		t.Fatalf("remove returned false on target %d", target)
	}

	if bstCount.Count != c-1 {
		t.Fatalf("unexpected count after remove - expecting %d, got %d", c, bstCount.Count)
	}

	if bstCount.Find(target) {
		t.Fatalf("found removed target %d", target)
	}

	for i := start; i < limit; i++ {
		if !bstCount.Find(i) {
			if i == target {
				continue
			}

			t.Fatalf("missing node %d", i)
		}
	}
}

func TestBstCountCustomInsertFind(t *testing.T) {
	bst := new(BstCmp[*big.Int])
	bstCount := NewBstCount[*big.Int](bst)

	start := int64(1)
	limit := int64(10)

	var c uint64
	for i := start; i < limit; i++ {
		inserted := bstCount.Insert(big.NewInt(i))
		if !inserted {
			t.Fatalf("Insert returned false after inserting new value %d", i)
		}

		c++
		if bstCount.Count != c {
			t.Fatalf("Unexpected BstCount count - expecting %d, got %d", c, bstCount.Count)
		}
	}

	for i := start; i < limit; i++ {
		inserted := bstCount.Insert(big.NewInt(i))
		if inserted {
			t.Fatalf("Insert returned true after replacing value %d with new node", i)
		}

		if !bst.Find(big.NewInt(i)) {
			t.Fatalf("missing node %d", i)
		}
	}

	if bstCount.Count != c {
		t.Fatalf("Unexpected BstCount count after reinsert duplicates - expecting %d, got %d", c, bstCount.Count)
	}

	outOfRange := int64(-1)
	if bst.Find(big.NewInt(outOfRange)) {
		t.Fatalf("unexpected false positive for %d", outOfRange)
	}

	curr := &bst.Root
	for curr != nil {
		t.Logf("node %+v", curr)
		curr = curr.right
	}
}
