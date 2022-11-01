package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyArrayEquals(t *testing.T) {
	emptyA := DependencyArray{}
	emptyB := DependencyArray{}

	assert.True(t, emptyA.Equals(emptyB))

	oneA := DependencyArray{
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
	}
	oneB := DependencyArray{
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
	}

	assert.True(t, oneA.Equals(oneB))

	twoA := DependencyArray{
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
		NewDependency("github.com/b/b", NewVersion(Tag, "beta")),
	}
	twoB := DependencyArray{
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
		NewDependency("github.com/b/b", NewVersion(Branch, "beta")),
	}

	assert.False(t, oneA.Equals(twoA))

	assert.False(t, twoA.Equals(twoB))
}
