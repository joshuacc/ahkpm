package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"os"
)

// Manifest contains the data from ahkpm.json
type Manifest struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	Repository   string `json:"repository"`
	Website      string `json:"website"`
	License      string `json:"license"`
	IssueTracker string `json:"issueTracker"`
	Author       Person `json:"author"`
	dependencies []Dependency
}

type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

func NewManifest() *Manifest {
	return &Manifest{
		Author:       Person{},
		dependencies: []Dependency{},
	}
}

func (m *Manifest) String() string {
	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		utils.Exit("Error marshalling ahkpm.json to string")
	}
	return string(jsonBytes)
}

func ManifestFromCwd() *Manifest {
	m, err := ManifestFromFile("ahkpm.json")
	if err != nil {
		utils.Exit(err.Error())
	}
	return m
}

func ManifestFromFile(path string) (*Manifest, error) {
	m := NewManifest()
	return utils.StructFromFile(path, m)
}

func (m Manifest) SaveToCwd() Manifest {
	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		utils.Exit("Error marshalling ahkpm.json to bytes")
	}
	err = os.WriteFile("ahkpm.json", jsonBytes, 0644)
	if err != nil {
		utils.Exit("Error writing ahkpm.json")
	}
	return m
}

func (m *Manifest) Dependencies() []Dependency {
	return m.dependencies
}

func (m *Manifest) AddDependency(name string, version Version) Manifest {
	newDep := NewDependency(name, version)

	foundIndex := -1
	for i, dep := range m.dependencies {
		if dep.Name() == name {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		m.dependencies = append(m.dependencies, newDep)
	} else {
		m.dependencies[foundIndex] = newDep
	}

	return *m
}

func (m *Manifest) UnmarshalJSON(data []byte) error {
	type Alias Manifest
	aux := &struct {
		Dependencies map[string]string `json:"dependencies"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for packageName, versionSpecifier := range aux.Dependencies {
		dep, err := DependencyFromSpecifiers(packageName, versionSpecifier)
		if err != nil {
			return err
		}
		m.dependencies = append(m.dependencies, dep)
	}

	return nil
}

func (m Manifest) MarshalJSON() ([]byte, error) {
	type Alias Manifest
	aux := &struct {
		*Alias
		Dependencies map[string]string `json:"dependencies"`
	}{
		Alias: (*Alias)(&m),
	}
	aux.Dependencies = make(map[string]string)
	for _, dep := range m.dependencies {
		aux.Dependencies[dep.Name()] = dep.Version().String()
	}

	return json.Marshal(aux)
}
