package core_test

import (
	. "ahkpm/src/core"
	"ahkpm/src/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveWithNoDependencies(t *testing.T) {
	dr := NewDependencyResolver()
	resolvedList, err := dr.Resolve([]Dependency{})

	assert.NoError(t, err)
	assert.Equal(t, []Node[Dependency]{}, resolvedList)
}

func TestResolveWithNoChildDependencies(t *testing.T) {
	dr := NewDependencyResolver()
	deps := []Dependency{
		NewDependency("foo", NewVersion(SemVerExact, "1.2.3")),
	}
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetPackageDependencies", deps[0]).Return([]Dependency{}, nil)

	dr.ReplacePackagesRepository(mockPR)

	resolvedList, err := dr.Resolve(deps)

	expectedList := []Node[Dependency]{
		{
			Value:    deps[0],
			Children: []Node[Dependency]{},
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedList, resolvedList)
}

func TestResolveWithChildDependencies(t *testing.T) {
	dr := NewDependencyResolver()
	deps := []Dependency{
		NewDependency("foo", NewVersion(SemVerExact, "1.2.3")),
	}
	childDeps := []Dependency{
		NewDependency("bar", NewVersion(SemVerExact, "1.2.3")),
	}
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetPackageDependencies", deps[0]).Return(childDeps, nil)
	mockPR.On("GetPackageDependencies", childDeps[0]).Return([]Dependency{}, nil)

	dr.ReplacePackagesRepository(mockPR)

	resolvedList, err := dr.Resolve(deps)

	expectedList := []Node[Dependency]{
		{
			Value: deps[0],
			Children: []Node[Dependency]{
				{
					Value:    childDeps[0],
					Children: []Node[Dependency]{},
				},
			},
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedList, resolvedList)
}
