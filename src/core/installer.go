package core

import (
	"ahkpm/src/utils"
	"os"
)

type Installer struct{}

func (i Installer) Install(deps DepsArray) {
	pr := NewPackagesRepository()

	m := ManifestFromCwd()
	if deps.Equals(m.Dependencies()) {
		lm := LockManifestFromCwd()
		os.RemoveAll("ahkpm-modules")
		for _, resolvedDep := range lm.Resolved {
			err := pr.CopyPackage(resolvedDep, resolvedDep.InstallPath)
			if err != nil {
				utils.Exit(err.Error())
			}
		}

		return
	}

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
