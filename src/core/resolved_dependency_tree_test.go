package core_test

import (
	. "ahkpm/src/core"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvedDependencyTreeFlatten(t *testing.T) {
	r := ResolvedDependencyTree{}
	r = append(r, TreeNode[ResolvedDependency]{
		Value: ResolvedDependency{
			Name:    "github.com/a/a",
			Version: "1.0.0",
		},
		Children: []TreeNode[ResolvedDependency]{
			{
				Value: ResolvedDependency{
					Name:    "github.com/b/b",
					Version: "1.0.0",
				},
				Children: []TreeNode[ResolvedDependency]{},
			},
		},
	})

	expected := []ResolvedDependency{
		{
			Name:    "github.com/a/a",
			Version: "1.0.0",
		},
		{
			Name:    "github.com/b/b",
			Version: "1.0.0",
		},
	}

	assert.Equal(t, expected, r.Flatten())
}

func TestResolvedDependencyTreeForEach(t *testing.T) {
	r := ResolvedDependencyTree{}
	r = append(r, TreeNode[ResolvedDependency]{
		Value: ResolvedDependency{
			Name:    "github.com/a/a",
			Version: "1.0.0",
		},
		Children: []TreeNode[ResolvedDependency]{
			{
				Value: ResolvedDependency{
					Name:    "github.com/b/b",
					Version: "1.0.0",
				},
				Children: []TreeNode[ResolvedDependency]{},
			},
		},
	})

	expected := []ResolvedDependency{
		{
			Name:    "github.com/a/a",
			Version: "1.0.0",
		},
		{
			Name:    "github.com/b/b",
			Version: "1.0.0",
		},
	}

	var actual []ResolvedDependency
	err := r.ForEach(func(n TreeNode[ResolvedDependency]) error {
		actual = append(actual, n.Value)
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestResolvedDependencyTreeMap(t *testing.T) {
	r := ResolvedDependencyTree{}
	r = append(r, TreeNode[ResolvedDependency]{
		Value: ResolvedDependency{
			Name:    "github.com/a/a",
			Version: "1.0.1",
		},
		Children: []TreeNode[ResolvedDependency]{
			{
				Value: ResolvedDependency{
					Name:    "github.com/b/b",
					Version: "0.1.0",
				},
				Children: []TreeNode[ResolvedDependency]{},
			},
		},
	})

	expected := ResolvedDependencyTree{}
	expected = append(expected, TreeNode[ResolvedDependency]{
		Value: ResolvedDependency{
			Name:    "github.com/a/a",
			Version: "2.0.2",
		},
		Children: []TreeNode[ResolvedDependency]{
			{
				Value: ResolvedDependency{
					Name:    "github.com/b/b",
					Version: "0.2.0",
				},
				Children: []TreeNode[ResolvedDependency]{},
			},
		},
	})

	actual := r.Map(func(n TreeNode[ResolvedDependency]) TreeNode[ResolvedDependency] {
		n.Value.Version = strings.Replace(n.Value.Version, "1", "2", 2)
		return n
	})

	assert.Equal(t, expected, actual)
}
