package core

type TreeNode[T any] struct {
	Value    T
	Children []TreeNode[T]
	Parent   *TreeNode[T]
}

func NewTreeNode[T any](value T) TreeNode[T] {
	return TreeNode[T]{
		Parent:   nil,
		Value:    value,
		Children: []TreeNode[T]{},
	}
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

func (n TreeNode[T]) ForEach(callback func(n TreeNode[T]) error) error {
	if &n.Value != (*T)(nil) {
		err := callback(n)
		if err != nil {
			return err
		}
	}

	for _, child := range n.Children {
		err := child.ForEach(callback)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n TreeNode[T]) WithChildren(children []TreeNode[T]) TreeNode[T] {
	for i := range children {
		children[i].Parent = &n
	}

	n.Children = append(n.Children, children...)

	return n
}

func (n TreeNode[T]) WithParent(parent TreeNode[T]) TreeNode[T] {
	n.Parent = &parent
	return n
}
