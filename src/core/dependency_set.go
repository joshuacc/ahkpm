package core

import "encoding/json"

type DependencySet struct {
	_set map[string]Dependency
}

func NewDependencySet() DependencySet {
	return DependencySet{
		_set: make(map[string]Dependency),
	}
}

// Equals returns true if the two sets contain the same dependencies with the same versions
func (set DependencySet) Equals(other DependencySet) bool {
	if set.Len() != other.Len() {
		return false
	}

	for _, otherDep := range other.AsMap() {
		dep, ok := set._set[otherDep.Name()]
		if !ok || !dep.Equals(otherDep) {
			return false
		}
	}

	return true
}

// Len is the number of dependencies in the set
func (ds DependencySet) Len() int {
	return len(ds._set)
}

// AsArray returns the dependencies as an array
func (ds DependencySet) AsArray() []Dependency {
	deps := make([]Dependency, len(ds._set))
	i := 0
	for _, dep := range ds._set {
		deps[i] = dep
		i++
	}
	return deps
}

// AsMap returns the dependencies as a map of name to dependency
func (ds DependencySet) AsMap() map[string]Dependency {
	return ds._set
}

func (ds DependencySet) MarshalJSON() ([]byte, error) {
	nameAndVersion := make(map[string]string)
	for _, dep := range ds._set {
		nameAndVersion[dep.Name()] = dep.Version().String()
	}

	return json.Marshal(nameAndVersion)
}

func (ds *DependencySet) UnmarshalJSON(data []byte) error {
	nameAndVersion := make(map[string]string)
	if err := json.Unmarshal(data, &nameAndVersion); err != nil {
		return err
	}

	ds._set = make(map[string]Dependency)
	for name, version := range nameAndVersion {
		dep, err := DependencyFromSpecifiers(name, version)
		if err != nil {
			return err
		}
		ds._set[name] = dep
	}

	return nil
}

// AddDependency adds a dependency to the array, replacing any existing
func (ds DependencySet) AddDependency(newDep Dependency) DependencySet {
	ds._set[newDep.Name()] = newDep
	return ds
}
