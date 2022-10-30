package core

import (
	"ahkpm/src/utils"
	"os"
)

type Installer struct{}

func (i Installer) Install(deps []Dependency) {
	pr := NewPackagesRepository()
	resolver := NewDependencyResolver()
	resolvedDepTree, err := resolver.Resolve(deps)
	if err != nil {
		utils.Exit(err.Error())
	}

	os.RemoveAll("ahkpm-modules")
	for _, topDepNode := range resolvedDepTree {
		err = topDepNode.ForEach(func(resolvedDepNode TreeNode[ResolvedDependency]) error {
			resolvedDep := resolvedDepNode.Value
			return pr.CopyPackage(resolvedDep, resolvedDep.InstallPath)
		})
		if err != nil {
			utils.Exit(err.Error())
		}
	}

	NewLockManifest().
		WithDependencies(deps).
		WithResolved(resolvedDepTree).
		SaveToCwd()
}
