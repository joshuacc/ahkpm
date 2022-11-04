package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencySetAddDependency(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta")))

	assert.Equal(t, 1, ds.Len())
	assert.Equal(t, "github.com/a/a", ds.AsArray()[0].Name())
	assert.Equal(t, NewVersion(Tag, "beta"), ds.AsArray()[0].Version())
}

func TestDependencySetAddDependencyWithExistingDependency(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta")))

	expected := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta")))

	assert.Equal(t, 1, ds.Len())
	assert.Equal(t, expected, ds)
}

func TestDependencySetEquals(t *testing.T) {
	emptyA := NewDependencySet()
	emptyB := NewDependencySet()

	assert.True(t, emptyA.Equals(emptyB))

	oneA := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta")))
	oneB := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta")))

	assert.True(t, oneA.Equals(oneB))

	twoA := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Tag, "beta")))
	twoB := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Branch, "beta")))

	assert.False(t, oneA.Equals(twoA))

	assert.False(t, twoA.Equals(twoB))
}

func TestDependencySetAsArray(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Tag, "beta")))

	expected := []Dependency{
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
		NewDependency("github.com/b/b", NewVersion(Tag, "beta")),
	}
	assert.Equal(t, 2, len(ds.AsArray()))
	assert.Equal(t, expected, ds.AsArray())
}

func TestDependencySetLen(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Tag, "beta")))

	assert.Equal(t, 2, ds.Len())
}

func TestDependencySetAsMap(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Tag, "beta")))

	assert.Equal(t, 2, len(ds.AsMap()))
	assert.Equal(
		t,
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
		ds.AsMap()["github.com/a/a"],
	)
	assert.Equal(
		t,
		NewDependency("github.com/b/b", NewVersion(Tag, "beta")),
		ds.AsMap()["github.com/b/b"],
	)
}

func TestDependencySetMarshalJSON(t *testing.T) {
	ds := NewDependencySet().
		AddDependency(NewDependency("github.com/a/a", NewVersion(Tag, "beta"))).
		AddDependency(NewDependency("github.com/b/b", NewVersion(Branch, "beta")))

	expected := `{"github.com/a/a":"tag:beta","github.com/b/b":"branch:beta"}`

	actual, err := ds.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestDependencySetUnmarshalJSON(t *testing.T) {
	json := `{"github.com/a/a":"tag:beta","github.com/b/b":"branch:beta"}`

	ds := NewDependencySet()

	err := ds.UnmarshalJSON([]byte(json))

	assert.Nil(t, err)
	assert.Equal(t, 2, ds.Len())
	assert.Equal(
		t,
		NewDependency("github.com/a/a", NewVersion(Tag, "beta")),
		ds.AsMap()["github.com/a/a"],
	)
	assert.Equal(
		t,
		NewDependency("github.com/b/b", NewVersion(Branch, "beta")),
		ds.AsMap()["github.com/b/b"],
	)
}
