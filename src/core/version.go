package core

import (
	"ahkpm/src/utils"
	"errors"
	"strings"
)

type Version struct {
	Kind  VersionKind
	Value string
}

type VersionKind string

const (
	SemVerExact VersionKind = "Semantic Version"
	Branch      VersionKind = "Branch"
	Tag         VersionKind = "Tag"
	Commit      VersionKind = "Commit"
)

func (v Version) FromString(versionSpecifier string) (Version, error) {
	if utils.IsSemVer(versionSpecifier) {
		v.Kind = SemVerExact
		v.Value = versionSpecifier
	} else if strings.HasPrefix(versionSpecifier, "branch:") {
		v.Kind = Branch
		v.Value = strings.TrimPrefix(versionSpecifier, "branch:")
	} else if strings.HasPrefix(versionSpecifier, "tag:") {
		v.Kind = Tag
		v.Value = strings.TrimPrefix(versionSpecifier, "tag:")
	} else if strings.HasPrefix(versionSpecifier, "commit:") {
		v.Kind = Commit
		v.Value = strings.TrimPrefix(versionSpecifier, "commit:")
	} else {
		return v, errors.New("Invalid version specifier " + versionSpecifier)
	}

	return v, nil
}

func (v Version) String() string {
	if v.Kind == Branch || v.Kind == Tag || v.Kind == Commit {
		return strings.ToLower(string(v.Kind)) + ":" + v.Value
	}
	if v.Kind == SemVerExact {
		return v.Value
	}
	utils.Exit("Invalid version kind " + string(v.Kind))
	return ""
}
