package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenWithoutChildren(t *testing.T) {
	n := TreeNode[string]{
		Value:    "foo",
		Children: []TreeNode[string]{},
	}

	assert.Equal(t, []string{"foo"}, n.Flatten())
}

func TestFlattenWithChildren(t *testing.T) {
	n := TreeNode[string]{
		Value: "foo",
		Children: []TreeNode[string]{
			{
				Value:    "bar",
				Children: []TreeNode[string]{},
			},
			{
				Value: "baz",
				Children: []TreeNode[string]{
					{
						Value:    "qux",
						Children: []TreeNode[string]{},
					},
				},
			},
		},
	}

	assert.Equal(t, []string{"foo", "bar", "baz", "qux"}, n.Flatten())
}

func TestWithChildren(t *testing.T) {
	n := TreeNode[string]{
		Value:    "foo",
		Children: []TreeNode[string]{},
	}

	n = n.WithChildren([]TreeNode[string]{
		{
			Value:    "bar",
			Children: []TreeNode[string]{},
		},
	})

	parent := TreeNode[string]{
		Value:    "foo",
		Children: []TreeNode[string]{},
	}
	child := TreeNode[string]{
		Value:    "bar",
		Children: []TreeNode[string]{},
		Parent:   &parent,
	}
	parent.Children = append(parent.Children, child)

	assert.Equal(t, parent, n)
}
