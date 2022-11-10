package cmd_test

import (
	"ahkpm/src/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUpdatedVersion(t *testing.T) {
	type Case struct {
		currentVersion string
		newVersionSpec string
		expected       string
		shouldError    bool
	}

	cases := []Case{
		{"1.2.3", "major", "2.0.0", false},
		{"1.2.3", "minor", "1.3.0", false},
		{"1.2.3", "patch", "1.2.4", false},
		{"1.2.3", "9.8.7", "9.8.7", false},
		{"1.2.3", "foobar", "", true},
	}

	for _, c := range cases {
		v, err := cmd.GetUpdatedVersion(c.currentVersion, c.newVersionSpec)
		if c.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.expected, v)
		}
	}
}
