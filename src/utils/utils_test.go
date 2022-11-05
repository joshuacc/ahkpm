package utils_test

import (
	. "ahkpm/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSemver(t *testing.T) {
	assert.True(t, IsSemVer("1.2.3"))
	assert.True(t, IsSemVer("1.2.3-beta.1"))
	assert.True(t, IsSemVer("1.2.3-beta.1+build.1"))
	assert.False(t, IsSemVer("foobar"))
}
