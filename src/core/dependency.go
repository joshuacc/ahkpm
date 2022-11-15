package core

import (
	"ahkpm/src/invariant"
	. "ahkpm/src/service_locator"
	"errors"
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
func DependencyFromSpecifiers(name string, versionSpecifier string, maybeLocator ...*ServiceLocator) (Dependency, error) {
	name = CanonicalizeDependencyName(name)

	isValidName := isValidDependencyName(name)
	if !isValidName {
		return nil, errors.New("Invalid dependency name " + name)
	}

	dep := dependency{name: name}

	if versionSpecifier == "" {
		locator := GetServiceLocator(maybeLocator)
		pr := locator.Get("PackagesRepository").(PackagesRepository)
		latestVersion, err := pr.GetLatestVersion(name)
		if err != nil {
			return nil, err
		}
		dep.version = latestVersion
		// If we got back a semantic version, convert it to a range so that we
		// will get the latest version on `ahkpm update` in the future
		if latestVersion.Kind() == SemVerExact {
			dep.version = NewVersion(SemVerRange, "^"+latestVersion.Value())
		}
	} else {
		version, err := VersionFromSpecifier(versionSpecifier)
		if err != nil {
			return nil, err
		}
		dep.version = version
	}

	return dep, nil
}

func DependencyFromSpecifier(specifier string) (Dependency, error) {
	splitSpecifier := strings.SplitN(specifier, "@", 2)
	packageName := splitSpecifier[0]
	versionSpecifier := ""
	if len(splitSpecifier) == 2 {
		versionSpecifier = splitSpecifier[1]
	}

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
	invariant.AssertNoError(err)

	return isMatch
}

func CanonicalizeDependencyName(name string) string {
	if strings.HasPrefix(name, "gh:") {
		name = strings.Replace(name, "gh:", "github.com/", 1)
	}

	return name
}
