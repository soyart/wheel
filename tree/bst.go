package tree

import (
	"golang.org/x/exp/constraints"
)

// Bst is BST implementation with nodeWrapper as node.
type Bst[T constraints.Ordered] struct {
	Root BinaryTreeNodeWrapper[T]
}

func NewBst[T constraints.Ordered]() *Bst[T] {
	return new(Bst[T])
}

func (b *Bst[T]) Insert(item T) bool {
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

	return BstInsert(
		&b.Root,
		&node,
	)
}

func (b *Bst[T]) Remove(target T) bool {
	return BstRemove(&b.Root, target) != nil
}

func (b *Bst[T]) Find(target T) bool {
	return BstFind(&b.Root, target)
}

func BstInsert[T constraints.Ordered](root *BinaryTreeNodeWrapper[T], node *BinaryTreeNodeWrapper[T]) bool {
	curr := root

	for {
		switch {
		// Found leaf node
		case curr == nil:
			*curr = *node

			return true

		// Do nothing if duplicate nodes
		case node.value == curr.value:
			exists := curr.ok

			node.ok = true
			curr = node

			return !exists

		case node.value < curr.value:
			if curr.left == nil {
				curr.left = node
				return true
			}

			curr = curr.left

		case node.value > curr.value:
			if curr.right == nil {
				curr.right = node
				return true
			}

			curr = curr.right
		}
	}
}

func BstFind[T constraints.Ordered](root *BinaryTreeNodeWrapper[T], target T) bool {
	curr := root

	for {
		if curr == nil || curr.IsNull() {
			return false
		}

		switch {
		case target == curr.value:
			return true

		case target < curr.value:
			curr = curr.left

		case target > curr.value:
			curr = curr.right

		default:
			panic("unhandled case")
		}
	}
}

// BstRemove removes target from subtree tree, returning the new root of the subtree
func BstRemove[T constraints.Ordered](root *BinaryTreeNodeWrapper[T], target T) *BinaryTreeNodeWrapper[T] {
	switch {
	case root == nil:
		return nil

	case target < root.value:
		root.left = BstRemove(root.left, target)

	case target > root.value:
		root.right = BstRemove(root.right, target)

	// Do the actual removal
	default:
		switch {
		case root.left == nil:
			return root.right

		case root.right == nil:
			return root.left

		default:
			replacement := digLeft(root.right)

			root.ok = false
			root.value = replacement.value
			root.right = BstRemove(root.right, replacement.value)
		}
	}

	return root
}

// Returns left leaf of a tree root
func digLeft[T any](root *BinaryTreeNodeWrapper[T]) *BinaryTreeNodeWrapper[T] {
	curr := root
	for curr.left != nil && curr.left.ok {
		curr = curr.left
	}

	return curr
}

// Returns right leaf of a tree root
//
//nolint:unused
func digRight[T any](root *BinaryTreeNodeWrapper[T]) *BinaryTreeNodeWrapper[T] {
	curr := root
	for curr.right != nil && curr.right.ok {
		curr = curr.right
	}

	return curr
}

func BstInsertRecurse[T constraints.Ordered](root *BinaryTreeNodeWrapper[T], node *BinaryTreeNodeWrapper[T]) bool {
	switch {
	case !root.ok:
		root = node
		root.ok = true

		return true

	case node.value == root.value:
		exists := root.ok

		node.ok = true
		root = node

		return !exists

	case node.value < root.value:
		if root.left == nil {
			root.left = node
			return true
		}

		return BstInsertRecurse(root.left, node)

	case node.value > root.value:
		if root.right == nil {
			root.right = node
			return true
		}

		return BstInsertRecurse(root.right, node)
	}

	return false
}
