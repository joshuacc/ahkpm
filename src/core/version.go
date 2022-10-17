package core

import (
	"ahkpm/src/utils"
	"strings"
)

type Version struct {
	Kind  VersionKind
	Value string
}

type VersionKind string

const (
	SemVerExact VersionKind = "SemVerExact"
	SemVerRange VersionKind = "SemVerRange"
	Branch      VersionKind = "Branch"
	Tag         VersionKind = "Tag"
	Commit      VersionKind = "Commit"
)

func (v Version) FromString(versionSpecifier string) Version {
	if utils.IsSemVer(versionSpecifier) {
		v.Kind = SemVerExact
		v.Value = versionSpecifier
	} else if utils.IsSemVerRange(versionSpecifier) {
		v.Kind = SemVerRange
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
		utils.Exit("Invalid version string " + versionSpecifier)
	}

	return v
}

func (v Version) String() string {
	if v.Kind == Branch || v.Kind == Tag || v.Kind == Commit {
		return strings.ToLower(string(v.Kind)) + ":" + v.Value
	}
	utils.Exit("Invalid version kind " + string(v.Kind))
	return ""
}
