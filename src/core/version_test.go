package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionFromString(t *testing.T) {
	type Case struct {
		specifier   string
		kind        VersionKind
		value       string
		shouldError bool
	}

	cases := [5]Case{
		{"1.2.3", SemVerExact, "1.2.3", false},
		{"branch:master", Branch, "master", false},
		{"tag:1.2.3", Tag, "1.2.3", false},
		{"commit:1234567890", Commit, "1234567890", false},
		{"foobar", SemVerExact, "", true},
	}

	for _, c := range cases {
		v, err := VersionFromSpecifier(c.specifier)

		if c.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.kind, v.Kind())
			assert.Equal(t, c.value, v.Value())
		}
	}
}

func TestVersionString(t *testing.T) {
	type Case struct {
		kind     VersionKind
		value    string
		expected string
	}

	cases := [4]Case{
		{SemVerExact, "1.2.3", "1.2.3"},
		{Branch, "master", "branch:master"},
		{Tag, "beta", "tag:beta"},
		{Commit, "1234567890", "commit:1234567890"},
	}

	for _, c := range cases {
		v := NewVersion(c.kind, c.value)
		assert.Equal(t, c.expected, v.String())
	}
}
