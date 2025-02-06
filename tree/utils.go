package tree

func LeftChildIdx(parent int) int {
	return (2 * parent) + 1
}

func RightChildIdx(parent int) int {
	return LeftChildIdx(parent) + 1
}

func ParentIdx(child int) int {
	rightChild := child%2 == 0
	if rightChild {
		return (child - 2) / 2
	}

	return (child - 1) / 2
}
