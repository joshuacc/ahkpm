package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"errors"
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
	dependencies map[string]string
}

type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

func NewManifest() *Manifest {
	return &Manifest{
		Author:       Person{},
		dependencies: make(map[string]string),
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
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error reading ahkpm.json at " + path)
	}
	m := NewManifest()
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return nil, errors.New("Error unmarshalling ahkpm.json")
	}
	return m, nil
}

func (m Manifest) Save() Manifest {
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
	deps := make([]Dependency, 0)
	for name, versionString := range m.dependencies {
		version, err := VersionFromSpecifier(versionString)
		if err != nil {
			utils.Exit("Error parsing version specifier " + versionString)
		}
		deps = append(deps, NewDependency(name, version))
	}
	return deps
}

func (m *Manifest) AddDependency(name string, version Version) Manifest {
	m.dependencies[name] = version.String()
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
	m.dependencies = aux.Dependencies
	return nil
}

func (m Manifest) MarshalJSON() ([]byte, error) {
	type Alias Manifest
	foo := &struct {
		*Alias
		Dependencies map[string]string `json:"dependencies"`
	}{
		Alias:        (*Alias)(&m),
		Dependencies: m.dependencies,
	}

	return json.Marshal(foo)
}
