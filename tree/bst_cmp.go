package tree

import (
	"github.com/soyart/wheel"
)

// BstCmp is for types with Cmp(T) int methods, e.g. *big.Int
type BstCmp[T wheel.CmpOrdered[T]] struct {
	Root BinaryTreeNodeWrapper[T]
}

func NewBstCmp[T wheel.CmpOrdered[T]]() *BstCmp[T] {
	return &BstCmp[T]{}
}

func (b *BstCmp[T]) Insert(item T) bool {
	node := BinaryTreeNodeWrapper[T]{
		value: item,
		ok:    true,
		left:  nil,
		right: nil,
	}

	if !b.Root.ok {
		b.Root = node
		return true
	}

	return BstCmpInsert(
		&b.Root,
		&node,
	)
}

func (b *BstCmp[T]) Find(target T) bool {
	return BstCmpFind(&b.Root, target)
}

func (b *BstCmp[T]) Remove(target T) bool {
	return BstCmpRemove(&b.Root, target) != nil
}

func BstCmpInsert[T wheel.CmpOrdered[T]](root, node *BinaryTreeNodeWrapper[T]) bool {
	curr := root

	for {
		if curr == nil {
			*curr = *node

			return true
		}

		cmp := node.Value().Cmp(curr.Value())

		switch {
		case cmp == 0:
			exists := curr.ok
			node.ok = true
			curr = node

			return !exists

		case cmp < 0:
			if curr.left == nil {
				curr.left = node
				return true
			}

			curr = curr.left

		case cmp > 0:
			if curr.right == nil {
				curr.right = node
				return true
			}

			curr = curr.right
		}
	}
}

func BstCmpFind[T wheel.CmpOrdered[T]](root *BinaryTreeNodeWrapper[T], target T) bool {
	curr := root

	for {
		if curr == nil || curr.IsNull() {
			return false
		}

		cmp := target.Cmp(curr.value)

		switch {
		case cmp == 0:
			return true

		case cmp < 0:
			curr = curr.left

		case cmp > 0:
			curr = curr.right

		default:
			panic("unhandled case")
		}
	}
}

// BstCmpRemove removes target from subtree tree, returning the new root of the subtree
func BstCmpRemove[T wheel.CmpOrdered[T]](root *BinaryTreeNodeWrapper[T], target T) *BinaryTreeNodeWrapper[T] {
	if root == nil {
		return nil
	}

	cmp := target.Cmp(root.value)

	switch {
	case cmp < 0:
		root.left = BstCmpRemove(root.left, target)

	case cmp > 0:
		root.right = BstCmpRemove(root.right, target)

	// Do the actual removal
	case cmp == 0:
		switch {
		case root.left == nil:
			return root.right

		case root.right == nil:
			return root.left

		default:
			replacement := digLeft(root.right)

			root.ok = false
			root.value = replacement.value
			root.right = BstCmpRemove(root.right, replacement.value)
		}
	}

	return root
}
