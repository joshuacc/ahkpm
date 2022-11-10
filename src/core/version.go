package core

import (
	"ahkpm/src/invariant"
	"ahkpm/src/utils"
	"errors"
	"regexp"
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
	SemVerRange VersionKind = "Semantic Version Range"
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
	} else if utils.IsSemVerRange(versionSpecifier) {
		v.kind = SemVerRange
		v.value = getLegibleRange(versionSpecifier)
	} else {
		return v, errors.New("Invalid version specifier " + versionSpecifier)
	}

	return v, nil
}

// Check to see if the range is of form "1" or "1.2" to convert them
// to equivalents "1.x.x" and "1.2.x" respectively. The goal is to make
// the range more explicit and readable to the user before saving them.
func getLegibleRange(versionSpecifier string) string {
	isSimpleRange, err := regexp.Match(`^\d+\.?(\d+)?$`, []byte(versionSpecifier))
	invariant.AssertNoError(err)
	if isSimpleRange {
		if strings.Contains(versionSpecifier, ".") {
			// If it's of form "1.2", convert it to "1.2.x"
			return versionSpecifier + ".x"
		} else {
			// If it's of form "1", convert it to "1.x.x"
			return versionSpecifier + ".x.x"
		}
	}
	return versionSpecifier
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
