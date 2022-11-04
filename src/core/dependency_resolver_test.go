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
	resolvedList, err := dr.Resolve(NewDependencySet())

	assert.NoError(t, err)
	assert.Equal(t, []TreeNode[ResolvedDependency]{}, resolvedList)
}

func TestResolveWithNoChildDependencies(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)
	dep1 := NewDependency("github.com/ahkpm/ahkpm", NewVersion(SemVerExact, "1.2.3"))
	deps := NewDependencySet().AddDependency(dep1)
	partiallyResolvedDep := ResolvedDependency{
		Name:    dep1.Name(),
		Version: dep1.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDep := partiallyResolvedDep.WithDependencies(NewDependencySet())
	fullyResolvedDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm"
	mockPR.On("GetResolvedDependencySHA", dep1).Return(partiallyResolvedDep.SHA, nil)
	newSet := NewDependencySet()
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return(&newSet, nil)

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
	dep1 := NewDependency("github.com/ahkpm/ahkpm", NewVersion(SemVerExact, "1.2.3"))
	deps := NewDependencySet().AddDependency(dep1)
	partiallyResolvedDep := ResolvedDependency{
		Name:    dep1.Name(),
		Version: dep1.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDep := partiallyResolvedDep.WithDependencies(NewDependencySet())
	fullyResolvedDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm"

	childDep1 := NewDependency("github.com/abcd/abcd", NewVersion(SemVerExact, "1.2.3"))
	childDeps := NewDependencySet().AddDependency(childDep1)

	partiallyResolvedChildDep := ResolvedDependency{
		Name:    childDep1.Name(),
		Version: childDep1.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDep := partiallyResolvedChildDep.WithDependencies(NewDependencySet())
	fullyResolvedChildDep.InstallPath = "ahkpm-modules/github.com/ahkpm/ahkpm/ahkpm-modules/github.com/abcd/abcd"

	mockPR.On("GetResolvedDependencySHA", dep1).Return(partiallyResolvedDep.SHA, nil)
	emptySet := NewDependencySet()
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDep).Return(&emptySet, nil)
	mockPR.On("GetResolvedDependencySHA", childDep1).Return(partiallyResolvedChildDep.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return(&childDeps, nil)

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
	depA := NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3"))
	depB := NewDependency("github.com/b/b", NewVersion(SemVerExact, "1.2.3"))
	deps := NewDependencySet().AddDependency(depA).AddDependency(depB)
	badDepA := NewDependency("github.com/bad/dep", NewVersion(SemVerExact, "1.2.3"))
	aDeps := NewDependencySet().AddDependency(badDepA)
	badDepB := NewDependency("github.com/bad/dep", NewVersion(SemVerExact, "1.2.4"))
	bDeps := NewDependencySet().AddDependency(badDepB)
	partiallyResolvedDepA := ResolvedDependency{
		Name:    depA.Name(),
		Version: depA.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDepA := partiallyResolvedDepA.WithDependencies(aDeps)
	fullyResolvedDepA.InstallPath = "ahkpm-modules/github.com/a/a"

	partiallyResolvedDepB := ResolvedDependency{
		Name:    depB.Name(),
		Version: depB.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedDepB := partiallyResolvedDepB.WithDependencies(bDeps)
	fullyResolvedDepB.InstallPath = "ahkpm-modules/github.com/b/b"

	partiallyResolvedChildDepA := ResolvedDependency{
		Name:    badDepA.Name(),
		Version: badDepA.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDepA := partiallyResolvedChildDepA.WithDependencies(NewDependencySet())
	fullyResolvedChildDepA.InstallPath = "ahkpm-modules/github.com/a/a/ahkpm-modules/github.com/bad/dep"

	partiallyResolvedChildDepB := ResolvedDependency{
		Name:    badDepB.Name(),
		Version: badDepB.Version().String(),
		SHA:     "1234567890",
	}
	fullyResolvedChildDepB := partiallyResolvedChildDepB.WithDependencies(NewDependencySet())
	fullyResolvedChildDepB.InstallPath = "ahkpm-modules/github.com/b/b/ahkpm-modules/github.com/bad/dep"

	mockPR.On("GetResolvedDependencySHA", depA).Return(partiallyResolvedDepA.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", depB).Return(partiallyResolvedDepB.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", badDepA).Return(partiallyResolvedChildDepA.SHA, nil)
	mockPR.On("GetResolvedDependencySHA", badDepB).Return(partiallyResolvedChildDepB.SHA, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDepA).Return(&aDeps, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedDepB).Return(&bDeps, nil)
	emptySet := NewDependencySet()
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDepA).Return(&emptySet, nil)
	mockPR.On("GetPackageDependencies", partiallyResolvedChildDepB).Return(&emptySet, nil)

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingDependencySHA(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)

	depA := NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3"))
	deps := NewDependencySet().AddDependency(depA)

	mockPR.On("GetResolvedDependencySHA", depA).Return("", errors.New("error getting SHA"))

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}

func TestResolveWithErrorGettingPackageDependencies(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	dr := NewDependencyResolver().WithPackagesRepository(mockPR)

	depA := NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.2.3"))
	deps := NewDependencySet().AddDependency(depA)
	partiallyResolvedDep := ResolvedDependency{
		Name:    depA.Name(),
		Version: depA.Version().String(),
		SHA:     "1234567890",
	}

	mockPR.On("GetResolvedDependencySHA", depA).Return(partiallyResolvedDep.SHA, nil)
	emptySet := NewDependencySet()
	mockPR.On("GetPackageDependencies", partiallyResolvedDep).Return(&emptySet, errors.New("error getting package dependencies"))

	_, err := dr.Resolve(deps)

	assert.Error(t, err)
}
