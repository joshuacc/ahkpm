package core_test

import (
	. "ahkpm/src/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManifest(t *testing.T) {
	m := NewManifest()
	assert.IsType(t, &Manifest{}, m)
	assert.Equal(t, Person{}, m.Author)
	assert.Equal(t, make(map[string]string), m.Dependencies())
}
