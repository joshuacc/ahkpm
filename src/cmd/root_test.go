package cmd_test

import (
	"ahkpm/src/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersions(t *testing.T) {
	output := cmd.GetVersions()
	// In a release build, the version is updated, but it is "development" in development builds.
	assert.Contains(t, output, "ahkpm: development")
}
