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

func TestResolvedDependencyTreeFromArray(t *testing.T) {
	arr := []ResolvedDependency{
		{
			Name:        "github.com/a/a",
			InstallPath: "ahkpm-modules/github.com/a/a",
		},
		{
			Name:        "github.com/b/b",
			InstallPath: "ahkpm-modules/github.com/b/b",
		},
		{
			Name:        "github.com/aa/aa",
			InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/aa/aa",
		},
		{
			Name:        "github.com/ba/ba",
			InstallPath: "ahkpm-modules/github.com/b/b/ahkpm-modules/github.com/ba/ba",
		},
		{
			Name:        "github.com/ab/ab",
			InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/ab/ab",
		},
		{
			Name:        "github.com/aba/aba",
			InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/ab/ab/ahkpm-modules/github.com/aba/aba",
		},
	}

	tree := ResolvedDependencyTreeFromArray(arr)

	assert.Equal(t, "github.com/a/a", tree[0].Value.Name)
	assert.Equal(t, "github.com/aa/aa", tree[0].Children[0].Value.Name)
	assert.Equal(t, "github.com/ab/ab", tree[0].Children[1].Value.Name)
	assert.Equal(t, "github.com/aba/aba", tree[0].Children[1].Children[0].Value.Name)
	assert.Equal(t, "github.com/b/b", tree[1].Value.Name)
	assert.Equal(t, "github.com/ba/ba", tree[1].Children[0].Value.Name)
}

// func TestFindByNamesPath(t *testing.T) {
// 	tree := ResolvedDependencyTree{
// 		{
// 			Value: ResolvedDependency{
// 				Name:        "github.com/a/a",
// 				InstallPath: "ahkpm-modules/github.com/a/a",
// 			},
// 			Children: []TreeNode[ResolvedDependency]{
// 				{
// 					Value: ResolvedDependency{
// 						Name:        "github.com/aa/aa",
// 						InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/aa/aa",
// 					},
// 				},
// 				{
// 					Value: ResolvedDependency{
// 						Name:        "github.com/ab/ab",
// 						InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/ab/ab",
// 					},
// 					Children: []TreeNode[ResolvedDependency]{
// 						{
// 							Value: ResolvedDependency{
// 								Name:        "github.com/aba/aba",
// 								InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/ab/ab/ahkpm-modules/github.com/aba/aba",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Value: ResolvedDependency{
// 				Name:        "github.com/b/b",
// 				InstallPath: "ahkpm-modules/github.com/b/b",
// 			},
// 		},
// 	}

// 	expected := NewTreeNode(ResolvedDependency{
// 		Name:        "github.com/aba/aba",
// 		InstallPath: "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/ab/ab/ahkpm-modules/github.com/aba/aba",
// 	})

// 	actual := FindByNamesPath(tree, []string{"github.com/a/a", "github.com/ab/ab", "github.com/aba/aba"})
// 	assert.Equal(t, expected.Value.Name, actual.Value.Name)
// }
