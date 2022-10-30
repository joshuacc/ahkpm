package core_test

import (
	"ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDependency(t *testing.T) {
	version := core.NewVersion("branch", "main")
	dep := core.NewDependency("github.com/ahkpm/ahkpm", version)
	assert.Equal(t, "github.com/ahkpm/ahkpm", dep.Name())
	assert.Equal(t, version, dep.Version())
}

func TestDependencyFromSpecifiers(t *testing.T) {
	dep, err := core.DependencyFromSpecifiers("github.com/ahkpm/ahkpm", "branch:main")
	assert.Nil(t, err)
	assert.Equal(t, "github.com/ahkpm/ahkpm", dep.Name())
	assert.Equal(t, "branch:main", dep.Version().String())
}

func TestDependencyFromSpecifiersWithInvalidInputs(t *testing.T) {
	dep1, err1 := core.DependencyFromSpecifiers("github.com/ahkpm/ahkpm", "invalid")
	assert.NotNil(t, err1)
	assert.Nil(t, dep1)

	dep2, err2 := core.DependencyFromSpecifiers("invalid", "branch:main")
	assert.NotNil(t, err2)
	assert.Nil(t, dep2)
}
