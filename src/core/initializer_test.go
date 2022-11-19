package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNewManifestWithDefaults(t *testing.T) {
	m := GetNewManifestWithDefaults()

	assert.IsType(t, &Manifest{}, m)
	assert.Equal(t, "0.0.1", m.Version)
	assert.Equal(t, "", m.Description)
	assert.Equal(t, "", m.Include)
	assert.Equal(t, "MIT", m.License)
	assert.IsType(t, Person{}, m.Author)
	assert.Equal(t, NewDependencySet(), m.Dependencies)
}
