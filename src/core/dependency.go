package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Dependency interface {
	Name() string
	Version() Version
	Equals(other Dependency) bool
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
	name = CanonicalizeDependencyName(name)

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

func DependencyFromSpecifier(specifier string) (Dependency, error) {
	if !strings.Contains(specifier, "@") {
		return nil, fmt.Errorf("Invalid dependency specifier: %s", specifier)
	}

	splitSpecifier := strings.SplitN(specifier, "@", 2)
	packageName := splitSpecifier[0]
	versionSpecifier := splitSpecifier[1]

	return DependencyFromSpecifiers(packageName, versionSpecifier)
}

func (d dependency) Name() string {
	return d.name
}

func (d dependency) Version() Version {
	return d.version
}

func (d dependency) Equals(other Dependency) bool {
	return d.name == other.Name() && d.version.Equals(other.Version())
}

func isValidDependencyName(name string) bool {
	isMatch, err := regexp.MatchString(`^github\.com\/[\w-\.]+\/[\w-\.]+$`, name)
	if err != nil {
		panic(err)
	}

	return isMatch
}

func CanonicalizeDependencyName(name string) string {
	if strings.HasPrefix(name, "gh:") {
		name = strings.Replace(name, "gh:", "github.com/", 1)
	}

	return name
}
