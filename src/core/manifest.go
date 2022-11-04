package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"os"
)

// Manifest contains the data from ahkpm.json
type Manifest struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Description  string        `json:"description"`
	Repository   string        `json:"repository"`
	Website      string        `json:"website"`
	License      string        `json:"license"`
	IssueTracker string        `json:"issueTracker"`
	Author       Person        `json:"author"`
	Dependencies DependencySet `json:"dependencies"`
}

type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

func NewManifest() *Manifest {
	return &Manifest{
		Author:       Person{},
		Dependencies: NewDependencySet(),
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
