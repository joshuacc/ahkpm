package core

import (
	"ahkpm/src/utils"
	"errors"
	"strings"
)

type Version interface {
	Kind() VersionKind
	Value() string
	String() string
	Equals(other Version) bool
}

type version struct {
	kind  VersionKind
	value string
}

type VersionKind string

const (
	SemVerExact VersionKind = "Semantic Version"
	Branch      VersionKind = "Branch"
	Tag         VersionKind = "Tag"
	Commit      VersionKind = "Commit"
)

// NewVersion creates a new version with the given kind and value. It does *not*
// check if the version is valid.
func NewVersion(kind VersionKind, value string) Version {
	return version{
		kind:  kind,
		value: value,
	}
}

// Converts a version specifier string into a Version.
func VersionFromSpecifier(versionSpecifier string) (Version, error) {
	v := version{}

	if utils.IsSemVer(versionSpecifier) {
		v.kind = SemVerExact
		v.value = versionSpecifier
	} else if strings.HasPrefix(versionSpecifier, "branch:") {
		v.kind = Branch
		v.value = strings.TrimPrefix(versionSpecifier, "branch:")
	} else if strings.HasPrefix(versionSpecifier, "tag:") {
		v.kind = Tag
		v.value = strings.TrimPrefix(versionSpecifier, "tag:")
	} else if strings.HasPrefix(versionSpecifier, "commit:") {
		v.kind = Commit
		v.value = strings.TrimPrefix(versionSpecifier, "commit:")
	} else {
		return v, errors.New("Invalid version specifier " + versionSpecifier)
	}

	return v, nil
}

// Represents the Version as a valid version specifier string.
func (v version) String() string {
	if v.kind == Branch || v.kind == Tag || v.kind == Commit {
		return strings.ToLower(string(v.kind)) + ":" + v.value
	}
	return v.value
}

func (v version) Kind() VersionKind {
	return v.kind
}

func (v version) Value() string {
	return v.value
}

func (v version) Equals(other Version) bool {
	return v.kind == other.Kind() && v.value == other.Value()
}
