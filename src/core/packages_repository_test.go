package core_test

import (
	. "ahkpm/src/core"
	"ahkpm/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearCache(t *testing.T) {
	clearedPath := ""
	pr := NewPackagesRepository().WithRemoveAll(func(path string) error {
		clearedPath = path
		return nil
	})
	err := pr.ClearCache()
	assert.Nil(t, err)
	assert.Equal(t, utils.GetAhkpmDir()+`\cache`, clearedPath)
}

func TestGetLatestVersionMatchingRangeFromArray(t *testing.T) {
	type Case struct {
		range_      string
		versions    []string
		expected    string
		shouldError bool
	}

	cases := []Case{
		{"1.2.3", []string{"1.2.3"}, "1.2.3", false},
		{"1.2.3", []string{"1.2.3", "1.2.4"}, "1.2.3", false},
		{"1", []string{"1.2.3", "1.2.4"}, "1.2.4", false},
		{"2", []string{"1.2.3", "1.2.4", "1.2.5"}, "", true},
		{"2.3.x", []string{"1.2.3", "2.1.0", "2.3.0", "2.3.1"}, "2.3.1", false},
	}

	for _, c := range cases {
		v, err := GetLatestVersionMatchingRangeFromArray(c.versions, c.range_)
		if c.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, c.expected, v)
		}
	}
}
