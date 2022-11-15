package core_test

import (
	. "ahkpm/src/core"
	"ahkpm/src/invariant"
	"ahkpm/src/mocks"
	"ahkpm/src/service_locator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDependency(t *testing.T) {
	version := NewVersion("branch", "main")
	dep := NewDependency("github.com/ahkpm/ahkpm", version)
	assert.Equal(t, "github.com/ahkpm/ahkpm", dep.Name())
	assert.Equal(t, version, dep.Version())
}

func TestDependencyFromSpecifiers(t *testing.T) {
	dep, err := DependencyFromSpecifiers("github.com/ahkpm/ahkpm", "branch:main")
	assert.Nil(t, err)
	assert.Equal(t, "github.com/ahkpm/ahkpm", dep.Name())
	assert.Equal(t, "branch:main", dep.Version().String())
}

func TestDependencyFromSpecifiersWithInvalidInputs(t *testing.T) {
	dep1, err1 := DependencyFromSpecifiers("github.com/ahkpm/ahkpm", "invalid")
	assert.NotNil(t, err1)
	assert.Nil(t, dep1)

	dep2, err2 := DependencyFromSpecifiers("invalid", "branch:main")
	assert.NotNil(t, err2)
	assert.Nil(t, dep2)
}

func TestDependencyFromSpecifiersWithEmptyVersion(t *testing.T) {
	mockPR := &mocks.MockPackagesRepository{}
	mockPR.On("GetLatestVersion", "github.com/a/a").Return(NewVersion(SemVerExact, "1.33.7"), nil)

	locator := service_locator.NewServiceLocator()
	err := locator.Add("PackagesRepository", mockPR)
	invariant.AssertNoError(err)

	dep1, err1 := DependencyFromSpecifiers("github.com/a/a", "", locator)

	expected := NewDependency("github.com/a/a", NewVersion(SemVerRange, "^1.33.7"))

	assert.Nil(t, err1)
	assert.Equal(t, expected, dep1)
}

func TestDependencyFromSpecifier(t *testing.T) {
	dep, err := DependencyFromSpecifier("github.com/ahkpm/ahkpm@branch:main")
	assert.Nil(t, err)
	assert.Equal(t, "github.com/ahkpm/ahkpm", dep.Name())
	assert.Equal(t, "branch:main", dep.Version().String())
}

func TestDependencyFromSpecifierWithInvalidInputs(t *testing.T) {
	dep1, err1 := DependencyFromSpecifier("github.com/ahkpm/ahkpm")
	assert.NotNil(t, err1)
	assert.Nil(t, dep1)
}

func TestEquals(t *testing.T) {
	cases := []struct {
		dep1    Dependency
		dep2    Dependency
		isEqual bool
	}{
		{
			dep1:    NewDependency("github.com/a/a", NewVersion("branch", "main")),
			dep2:    NewDependency("github.com/a/a", NewVersion("branch", "main")),
			isEqual: true,
		},
		{
			dep1:    NewDependency("github.com/a/a", NewVersion("branch", "main")),
			dep2:    NewDependency("github.com/a/a", NewVersion("branch", "dev")),
			isEqual: false,
		},
		{
			dep1:    NewDependency("github.com/a/a", NewVersion("branch", "main")),
			dep2:    NewDependency("github.com/b/b", NewVersion("branch", "main")),
			isEqual: false,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.isEqual, c.dep1.Equals(c.dep2))
	}
}

func TestCanonicalizeDependencyName(t *testing.T) {
	assert.Equal(t, "github.com/a/a", CanonicalizeDependencyName("gh:a/a"))
	assert.Equal(t, "github.com/a/a", CanonicalizeDependencyName("github.com/a/a"))
}
