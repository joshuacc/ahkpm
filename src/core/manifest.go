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

func (m *Manifest) ReadFromFile() *Manifest {
	jsonBytes, err := os.ReadFile("ahkpm.json")
	if err != nil {
		utils.Exit("Error reading ahkpm.json")
	}
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		utils.Exit("Error unmarshalling ahkpm.json")
	}
	return m
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

func (m *Manifest) Dependencies() map[string]string {
	return m.dependencies
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
