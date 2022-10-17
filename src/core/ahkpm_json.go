package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"fmt"
	"os"
)

type AhkpmJson struct {
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

func (aj AhkpmJson) New() AhkpmJson {
	aj.Author = Person{}
	aj.dependencies = make(map[string]string)
	return aj
}

func (aj AhkpmJson) String() string {
	jsonBytes, err := json.MarshalIndent(aj, "", "  ")
	if err != nil {
		utils.Exit("Error marshalling ahkpm.json to string")
	}
	return string(jsonBytes)
}

func (aj AhkpmJson) ReadFromFile() AhkpmJson {
	jsonBytes, err := os.ReadFile("ahkpm.json")
	if err != nil {
		utils.Exit("Error reading ahkpm.json")
	}
	err = json.Unmarshal(jsonBytes, &aj)
	if err != nil {
		utils.Exit("Error unmarshalling ahkpm.json")
	}
	return aj
}

func (aj AhkpmJson) Save() AhkpmJson {
	fmt.Println("Saving ahkpm.json")
	for k, v := range aj.dependencies {
		fmt.Printf("  %s: %s", k, v)
	}
	jsonBytes, err := json.MarshalIndent(aj, "", "  ")
	if err != nil {
		utils.Exit("Error marshalling ahkpm.json to bytes")
	}
	err = os.WriteFile("ahkpm.json", jsonBytes, 0644)
	if err != nil {
		utils.Exit("Error writing ahkpm.json")
	}
	return aj
}

func (aj AhkpmJson) Dependencies() map[string]string {
	return aj.dependencies
}

func (aj AhkpmJson) AddDependency(name string, version Version) AhkpmJson {
	aj.dependencies[name] = version.String()
	return aj
}

func (aj *AhkpmJson) UnmarshalJSON(data []byte) error {
	type Alias AhkpmJson
	aux := &struct {
		Dependencies map[string]string `json:"dependencies"`
		*Alias
	}{
		Alias: (*Alias)(aj),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	aj.dependencies = aux.Dependencies
	fmt.Println("Unmarshalled ahkpm.json")
	return nil
}

func (aj AhkpmJson) MarshalJSON() ([]byte, error) {
	fmt.Println("Marshalling ahkpm.json")
	for k, v := range aj.dependencies {
		fmt.Printf("  %s: %s", k, v)
	}
	type Alias AhkpmJson
	foo := &struct {
		*Alias
		Dependencies map[string]string `json:"dependencies"`
	}{
		Alias:        (*Alias)(&aj),
		Dependencies: aj.dependencies,
	}

	for k, v := range foo.Dependencies {
		fmt.Printf("  %s: %s", k, v)
	}
	return json.Marshal(foo)
}
