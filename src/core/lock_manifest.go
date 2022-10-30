package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"os"
)

type LockManifest struct {
	LockfileVersion string `json:"lockfileVersion"`
	dependencies    []Dependency
	Resolved        []ResolvedDependency `json:"resolved"`
}

func NewLockManifest() LockManifest {
	return LockManifest{
		LockfileVersion: "1",
		dependencies:    make([]Dependency, 0),
		Resolved:        make([]ResolvedDependency, 0),
	}
}

func (lm LockManifest) WithResolved(resDeps []TreeNode[ResolvedDependency]) LockManifest {
	resolved := make([]ResolvedDependency, 0)
	for _, depNode := range resDeps {
		err := depNode.ForEach(func(depNode TreeNode[ResolvedDependency]) error {
			resolved = append(resolved, depNode.Value)
			return nil
		})
		if err != nil {
			utils.Exit(err.Error())
		}
	}

	lm.Resolved = resolved

	return lm
}

func (lm LockManifest) WithDependencies(deps []Dependency) LockManifest {
	lm.dependencies = deps
	return lm
}

func (lm *LockManifest) Dependencies() []Dependency {
	return lm.dependencies
}

func (lm *LockManifest) AddDependency(name string, version Version) LockManifest {
	newDep := NewDependency(name, version)

	foundIndex := -1
	for i, dep := range lm.dependencies {
		if dep.Name() == name {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		lm.dependencies = append(lm.dependencies, newDep)
	} else {
		lm.dependencies[foundIndex] = newDep
	}

	return *lm
}

func (lm *LockManifest) UnmarshalJSON(data []byte) error {
	type Alias LockManifest
	aux := &struct {
		Dependencies map[string]string `json:"dependencies"`
		*Alias
	}{
		Alias: (*Alias)(lm),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for packageName, versionSpecifier := range aux.Dependencies {
		dep, err := DependencyFromSpecifiers(packageName, versionSpecifier)
		if err != nil {
			return err
		}
		lm.dependencies = append(lm.dependencies, dep)
	}

	return nil
}

func (lm LockManifest) MarshalJSON() ([]byte, error) {
	type Alias LockManifest
	aux := &struct {
		*Alias
		Dependencies map[string]string    `json:"dependencies"`
		Resolved     []ResolvedDependency `json:"resolved"` // Ensure resolved is last
	}{
		Alias: (*Alias)(&lm),
	}
	aux.Dependencies = make(map[string]string)
	for _, dep := range lm.dependencies {
		aux.Dependencies[dep.Name()] = dep.Version().String()
	}
	aux.Resolved = lm.Resolved

	return json.Marshal(aux)
}

func (lm LockManifest) SaveToCwd() LockManifest {
	jsonBytes, err := json.MarshalIndent(lm, "", "  ")
	if err != nil {
		utils.Exit("Error marshalling ahkpm.lock to bytes")
	}
	err = os.WriteFile("ahkpm.lock", jsonBytes, 0644)
	if err != nil {
		utils.Exit("Error writing ahkpm.lock")
	}
	return lm
}
