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
	assert.Equal(t, []TreeNode[ResolvedDependency]{}, resolvedList)
}

func TestResolveWithNoChildDependencies(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)
	deps := []Dependency{
		NewDependency("github.com/ahkpm/ahkpm", NewVersion(SemVerExact, "1.2.3")),
	}
	partiallyResolvedDep := ResolvedDependency{
		Name:    deps[0].Name(),
		Version: deps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDep := partiallyResolvedDep.WithDependencies([]Dependency{})
	fullyResolvedDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm"
	mockPR.On("GetResolvedDependencySHA", deps[0]).Return(partiallyResolvedDep.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return([]Dependency{}, nil)

	resolvedList, err := dr.Resolve(deps)

	expectedList := []TreeNode[ResolvedDependency]{
		NewTreeNode(fullyResolvedDep),
	}

	assert.NoError(t, err)
	assert.EqualValues(t, expectedList, resolvedList)
}

func TestResolveWithChildDependencies(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)
	deps := []Dependency{
		NewDependency("github.com/ahkpm/ahkpm", NewVersion(SemVerExact, "1.2.3")),
	}
	partiallyResolvedDep := ResolvedDependency{
		Name:    deps[0].Name(),
		Version: deps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDep := partiallyResolvedDep.WithDependencies([]Dependency{})
	fullyResolvedDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm"

	childDeps := []Dependency{
		NewDependency("github.com/abcd/abcd", NewVersion(SemVerExact, "1.2.3")),
	}

	partiallyResolvedChildDep := ResolvedDependency{
		Name:    childDeps[0].Name(),
		Version: childDeps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDep := partiallyResolvedChildDep.WithDependencies([]Dependency{})
	fullyResolvedChildDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm/ahkpm-modules/github.com/abcd/abcd"

	mockPR.On("GetResolvedDependencySHA", deps[0]).Return(partiallyResolvedDep.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDep).Return([]Dependency{}, nil)
	mockPR.On("GetResolvedDependencySHA", childDeps[0]).Return(partiallyResolvedChildDep.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return(childDeps, nil)

	resolvedList, err := dr.Resolve(deps)

	expectedList := []TreeNode[ResolvedDependency]{
		NewTreeNode(fullyResolvedDep.WithDependencies(childDeps)).
			WithChildren(
				[]TreeNode[ResolvedDependency]{NewTreeNode(fullyResolvedChildDep)},
			),
	}
	assert.NoError(t, err)
	assert.EqualValues(t, expectedList[0].Value, resolvedList[0].Value)
	assert.EqualValues(t, expectedList[0].Children[0].Value, resolvedList[0].Children[0].Value)
}

func TestResolveWithConflictingChildDepencyVersions(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)
	deps := []Dependency{
		NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3")),
		NewDependency("github.com/b/b", NewVersion(SemVerExact, "1.2.3")),
	}
	aDeps := []Dependency{
		NewDependency("github.com/bad/dep", NewVersion(SemVerExact, "1.2.3")),
	}
	bDeps := []Dependency{
		NewDependency("github.com/bad/dep", NewVersion(SemVerExact, "1.2.4")),
	}
	partiallyResolvedDepA := ResolvedDependency{
		Name:    deps[0].Name(),
		Version: deps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDepA := partiallyResolvedDepA.WithDependencies(aDeps)
	fullyResolvedDepA.InstallPath = "ahkpm-modules/github.com/a/a"

	partiallyResolvedDepB := ResolvedDependency{
		Name:    deps[1].Name(),
		Version: deps[1].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDepB := partiallyResolvedDepB.WithDependencies(bDeps)
	fullyResolvedDepB.InstallPath = "ahkpm-modules/github.com/b/b"

	partiallyResolvedChildDepA := ResolvedDependency{
		Name:    aDeps[0].Name(),
		Version: aDeps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDepA := partiallyResolvedChildDepA.WithDependencies([]Dependency{})
	fullyResolvedChildDepA.InstallPath = "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/bad/dep"

	partiallyResolvedChildDepB := ResolvedDependency{
		Name:    bDeps[0].Name(),
		Version: bDeps[0].Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDepB := partiallyResolvedChildDepB.WithDependencies([]Dependency{})
	fullyResolvedChildDepB.InstallPath = "ahkpm-modules/github.com/b/b/ahkpm-modules/github.com/bad/dep"

	mockPR.On("GetResolvedDependencySHA", deps[0]).Return(partiallyResolvedDepA.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", deps[1]).Return(partiallyResolvedDepB.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", aDeps[0]).Return(partiallyResolvedChildDepA.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", bDeps[0]).Return(partiallyResolvedChildDepB.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDepA).Return(aDeps, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDepB).Return(bDeps, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDepA).Return([]Dependency{}, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDepB).Return([]Dependency{}, nil)

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingDependencySHA(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)

	deps := []Dependency{
		NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3")),
	}

	mockPR.On("GetResolvedDependencySHA", deps[0]).Return("", errors.New("error getting SHA"))

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingPackageDependencies(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)

	deps := []Dependency{
		NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3")),
	}
	partiallyResolvedDep := ResolvedDependency{
		Name:    deps[0].Name(),
		Version: deps[0].Version().String(),
		SHA:     "1234567890",
	}

	mockPR.On("GetResolvedDependencySHA", deps[0]).Return(partiallyResolvedDep.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return([]Dependency{}, errors.New("error getting package dependencies"))

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}
