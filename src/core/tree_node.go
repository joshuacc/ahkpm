package core

type TreeNode[T any] struct {
	Value    T
	Children []TreeNode[T]
}

func (n TreeNode[T]) IsLeaf() bool {
	return len(n.Children) == 0
}

func (n TreeNode[T]) Flatten() []T {
	if n.IsLeaf() {
		return []T{n.Value}
	}

	flattened := []T{n.Value}
	for _, child := range n.Children {
		flattened = append(flattened, child.Flatten()...)
	}

	return flattened
}
