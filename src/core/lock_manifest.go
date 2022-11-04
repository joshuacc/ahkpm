package core

import (
	"ahkpm/src/utils"
	"encoding/json"
	"os"
)

type LockManifest struct {
	LockfileVersion string               `json:"lockfileVersion"`
	Dependencies    DependencySet        `json:"dependencies"`
	Resolved        []ResolvedDependency `json:"resolved"`
}

func NewLockManifest() LockManifest {
	return LockManifest{
		LockfileVersion: "1",
		Dependencies:    NewDependencySet(),
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

func (lm LockManifest) WithDependencies(deps DependencySet) LockManifest {
	lm.Dependencies = deps
	return lm
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

func LockManifestFromCwd() (*LockManifest, error) {
	lm, err := LockManifestFromFile("ahkpm.lock")
	if err != nil {
		return nil, err
	}
	return lm, nil
}

func LockManifestFromFile(path string) (*LockManifest, error) {
	lm := NewLockManifest()
	return utils.StructFromFile(path, &lm)
}
