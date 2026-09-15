package tree

import (
	"github.com/soyart/wheel/list"
)

// BinaryTreeBasic is basic, minimal binary tree with node type N
type BinaryTreeBasic[N any] interface {
	// Insert inserts a node to the tree,
	// returning bool indicating if a node was added.
	// False is returned if an existing node was replaced
	// with new node, leaving tree size unchanged.
	Insert(node N) bool

	// Remove removes node, returning whether the removal
	// was successful.
	Remove(node N) bool

	Find(node N) bool
}

// BinaryTree have extra methods to work with nodes.
// P is any type used for indexing a node,
// e.g. a bintree with backing arrays may use int as P.
type BinaryTree[P any, N any] interface {
	BinaryTreeBasic[N]

	Parent(node P) P
	LeftChild(node P) P
	RightChild(node P) P
	Node(pos P) N

	NodeIsRoot(node P) bool
	NodeIsNull(node P) bool
}

type BinaryTreeNode[T any] interface {
	Value() T
	Left() BinaryTreeNode[T]
	Right() BinaryTreeNode[T]
	IsNull() bool
}

type BinaryTreeNodeWrapper[T any] struct {
	value T
	ok    bool

	left  *BinaryTreeNodeWrapper[T]
	right *BinaryTreeNodeWrapper[T]
}

func (n *BinaryTreeNodeWrapper[T]) Left() BinaryTreeNode[T]  { return n.left }
func (n *BinaryTreeNodeWrapper[T]) Right() BinaryTreeNode[T] { return n.right }
func (n *BinaryTreeNodeWrapper[T]) Value() T                 { return n.value }

func (n *BinaryTreeNodeWrapper[T]) IsNull() bool {
	return !n.ok &&
		n.left == nil && n.right == nil
}

func (n *BinaryTreeNodeWrapper[T]) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

func Inorder[P any, N any](
	tree BinaryTree[P, N],
	node P,
	walk func(N) error,
) error {
	stack := list.NewStackSafe[P]()
	curr := node
	for !tree.NodeIsNull(curr) || !stack.IsEmpty() {
		for !tree.NodeIsNull(curr) {
			stack.Push(curr)
			curr = tree.LeftChild(curr)
		}
		curr = *stack.Pop()
		if err := walk(tree.Node(curr)); err != nil {
			return err
		}
		curr = tree.RightChild(curr)
	}

	return nil
}

func InorderRecurse[P any, N any](
	tree BinaryTree[P, N],
	node P,
	walk func(N) error,
) error {
	if err := InorderRecurse(tree, tree.LeftChild(node), walk); err != nil {
		return err
	}
	if err := walk(tree.Node(node)); err != nil {
		return err
	}
	return InorderRecurse(tree, tree.RightChild(node), walk)
}

func InorderNode[N BinaryTreeNode[any]](node N, f func(N) error) error {
	stack := list.NewStackSafe[N]()
	curr := node

	for !curr.IsNull() || !stack.IsEmpty() {
		for !curr.IsNull() {
			stack.Push(curr)
			curr = curr.Left().(N)
		}
		if err := f(curr); err != nil {
			return err
		}

		curr = *stack.Pop()
		curr = curr.Right().(N)
	}

	return nil
}

func InorderNodeRecurse[T any, N BinaryTreeNode[T]](node N, f func(N) error) error {
	if err := InorderNodeRecurse(node.Left().(N), f); err != nil {
		return err
	}
	if err := f(node); err != nil {
		return err
	}
	if err := InorderNodeRecurse(node.Right().(N), f); err != nil {
		return err
	}

	return nil
}

// DigRight digs for smallest values in the subtree, returning
// the node as well as the height to that node.
func DigRight[T any](node BinaryTreeNode[T]) (BinaryTreeNode[T], uint) {
	curr := node
	var height uint

	for curr.Right() != nil && !curr.Right().IsNull() {
		curr = curr.Right()
		height++
	}

	return curr, height
}

// DigLeft digs for smallest values in the subtree, returning
// the node as well as the height to that node.
func DigLeft[T any](root BinaryTreeNode[T]) (BinaryTreeNode[T], uint) {
	curr := root
	var height uint

	for curr.Left() != nil && !curr.Left().IsNull() {
		curr = curr.Left()
		height++
	}

	return curr, height
}
