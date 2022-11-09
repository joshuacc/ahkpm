package cmd_test

import (
	. "ahkpm/src/cmd"
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDependenciesForDisplay(t *testing.T) {
	set := NewDependencySet()
	set.AddDependency(NewDependency("github.com/a/a", NewVersion(SemVerExact, "1.0.0")))
	set.AddDependency(NewDependency("github.com/abc/abc", NewVersion(Branch, "main")))
	output := GetDependenciesForDisplay(set)

	expected := "Name              \tVersion\n"
	expected += "------------------\t-----------\n"
	expected += "github.com/a/a    \t1.0.0\n"
	expected += "github.com/abc/abc\tbranch:main\n"

	assert.Equal(t, expected, output)
}
