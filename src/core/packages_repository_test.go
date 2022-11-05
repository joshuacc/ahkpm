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
