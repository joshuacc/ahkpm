package core

import (
	"errors"
	"regexp"
)

type Dependency interface {
	Name() string
	Version() Version
}

type dependency struct {
	name    string
	version Version
}

// NewDependency creates a new dependency with the given name and version.
// It does *not* check if the dependency is valid.
func NewDependency(name string, version Version) Dependency {
	return dependency{
		name:    name,
		version: version,
	}
}

// DependencyFromSpecifiers creates a new dependency from the given name and
// version specifier. It performs validation on both.
func DependencyFromSpecifiers(name string, versionSpecifier string) (Dependency, error) {
	isValidName := isValidDependencyName(name)
	if !isValidName {
		return nil, errors.New("Invalid dependency name " + name)
	}
	version, err := VersionFromSpecifier(versionSpecifier)
	if err != nil {
		return nil, err
	}

	return dependency{
		name:    name,
		version: version,
	}, nil
}

func (d dependency) Name() string {
	return d.name
}

func (d dependency) Version() Version {
	return d.version
}

func isValidDependencyName(name string) bool {
	isMatch, err := regexp.MatchString("^github\\.com\\/[\\w-\\.]+\\/[\\w-\\.]+$", name)
	if err != nil {
		panic(err)
	}

	return isMatch
}
