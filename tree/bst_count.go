package tree

type BstCount[T any] struct {
	Count uint64

	BinaryTreeBasic[T]
}

func NewBstCount[T any](bst BinaryTreeBasic[T]) *BstCount[T] {
	return &BstCount[T]{
		Count:           0,
		BinaryTreeBasic: bst,
	}
}

func (b *BstCount[T]) Insert(item T) bool {
	if b.BinaryTreeBasic.Insert(item) {
		b.Count++
		return true
	}

	return false
}

func (b *BstCount[T]) Find(item T) bool {
	return b.BinaryTreeBasic.Find(item)
}

func (b *BstCount[T]) Remove(item T) bool {
	if b.BinaryTreeBasic.Remove(item) {
		b.Count--
		return true
	}

	return false
}
