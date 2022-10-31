package core

import "encoding/json"

type ResolvedDependency struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	SHA          string `json:"sha"`
	InstallPath  string `json:"installPath"`
	dependencies []Dependency
}

func (rd ResolvedDependency) Dependencies() []Dependency {
	return rd.dependencies
}

func (rd ResolvedDependency) WithDependencies(deps []Dependency) ResolvedDependency {
	rd.dependencies = deps
	return rd
}

func (rd *ResolvedDependency) AddDependency(name string, version Version) ResolvedDependency {
	newDep := NewDependency(name, version)

	foundIndex := -1
	for i, dep := range rd.dependencies {
		if dep.Name() == name {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		rd.dependencies = append(rd.dependencies, newDep)
	} else {
		rd.dependencies[foundIndex] = newDep
	}

	return *rd
}

func (ld *ResolvedDependency) UnmarshalJSON(data []byte) error {
	type Alias ResolvedDependency
	aux := &struct {
		Dependencies map[string]string `json:"dependencies"`
		*Alias
	}{
		Alias: (*Alias)(ld),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for packageName, versionSpecifier := range aux.Dependencies {
		dep, err := DependencyFromSpecifiers(packageName, versionSpecifier)
		if err != nil {
			return err
		}
		ld.dependencies = append(ld.dependencies, dep)
	}

	return nil
}

func (rd ResolvedDependency) MarshalJSON() ([]byte, error) {
	type Alias ResolvedDependency
	aux := &struct {
		*Alias
		Dependencies map[string]string `json:"dependencies"`
	}{
		Alias: (*Alias)(&rd),
	}
	aux.Dependencies = make(map[string]string)
	for _, dep := range rd.dependencies {
		aux.Dependencies[dep.Name()] = dep.Version().String()
	}

	return json.Marshal(aux)
}
