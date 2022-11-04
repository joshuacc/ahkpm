package core_test

import (
	. "ahkpm/src/core"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManifest(t *testing.T) {
	m := NewManifest()
	assert.IsType(t, &Manifest{}, m)
	assert.Equal(t, Person{}, m.Author)
	assert.Equal(t, NewDependencySet(), m.Dependencies)
}

func TestMarshalJSON(t *testing.T) {
	m := NewManifest()
	m.Name = "ahkpm"
	m.Version = "0.0.1"
	m.Description = "A package manager for AutoHotkey"
	m.Repository = "https://github.com/ahkpm/ahkpm"
	m.Website = "https://ahkpm.dev"
	m.License = "MIT"
	m.IssueTracker = "https://github.com/ahkpm/ahkpm/issues"
	m.Author = Person{
		Name:    "Thomas Aquinas",
		Email:   "angelicdoctor@example.com",
		Website: "https://en.wikipedia.org/wiki/Thomas_Aquinas",
	}
	dep := NewDependency("github.com/ahkpm/ahkpm", NewVersion("Branch", "main"))
	m.Dependencies.AddDependency(dep)
	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	assert.Nil(t, err)

	expected := `{
		"name": "ahkpm",
		"version": "0.0.1",
		"description": "A package manager for AutoHotkey",
		"repository": "https://github.com/ahkpm/ahkpm",
		"website": "https://ahkpm.dev",
		"license": "MIT",
		"issueTracker": "https://github.com/ahkpm/ahkpm/issues",
		"author": {
			"name": "Thomas Aquinas",
			"email": "angelicdoctor@example.com",
			"website": "https://en.wikipedia.org/wiki/Thomas_Aquinas"
		},
		"dependencies": {
			"github.com/ahkpm/ahkpm": "branch:main"
		}
	}`
	assert.JSONEq(t, expected, string(jsonBytes))
}
