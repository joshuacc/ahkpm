package core_test

import (
	. "ahkpm/src/core"
	"ahkpm/src/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveWithNoDependencies(t *testing.T) {
	dr := NewDependencyResolver()
	resolvedList, err := dr.Resolve([]Dependency{})

	assert.NoError(t, err)
	assert.Equal(t, []TreeNode[Dependency]{}, resolvedList)
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

	expectedList := []TreeNode[Dependency]{NewTreeNode(deps[0])}

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

	expectedList := []TreeNode[Dependency]{
		NewTreeNode(deps[0]).
			WithChildren(
				[]TreeNode[Dependency]{NewTreeNode(childDeps[0])},
			),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedList, resolvedList)
}

func TestResolveWithConflictingChildDepencyVersions(t *testing.T) {
	dr := NewDependencyResolver()
	deps := []Dependency{
		NewDependency("A", NewVersion(SemVerExact, "1.2.3")),
		NewDependency("B", NewVersion(SemVerExact, "1.2.3")),
	}
	aDeps := []Dependency{
		NewDependency("badDep", NewVersion(SemVerExact, "1.2.3")),
	}
	bDeps := []Dependency{
		NewDependency("badDep", NewVersion(SemVerExact, "1.2.4")),
	}
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetPackageDependencies", deps[0]).Return(aDeps, nil)
	mockPR.On("GetPackageDependencies", deps[1]).Return(bDeps, nil)
	mockPR.On("GetPackageDependencies", aDeps[0]).Return([]Dependency{}, nil)
	mockPR.On("GetPackageDependencies", bDeps[0]).Return([]Dependency{}, nil)

	dr.ReplacePackagesRepository(mockPR)

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingPackageDependencies(t *testing.T) {
	dr := NewDependencyResolver()
	deps := []Dependency{
		NewDependency("foo", NewVersion(SemVerExact, "1.2.3")),
	}
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetPackageDependencies", deps[0]).Return([]Dependency{}, errors.New("error getting package dependencies"))

	dr.ReplacePackagesRepository(mockPR)

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingPackageDependenciesForChildDependency(t *testing.T) {
	dr := NewDependencyResolver()
	deps := []Dependency{
		NewDependency("foo", NewVersion(SemVerExact, "1.2.3")),
	}
	childDeps := []Dependency{
		NewDependency("bar", NewVersion(SemVerExact, "1.2.3")),
	}
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetPackageDependencies", deps[0]).Return(childDeps, nil)
	mockPR.On("GetPackageDependencies", childDeps[0]).Return([]Dependency{}, errors.New("error getting package dependencies"))

	dr.ReplacePackagesRepository(mockPR)

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}
