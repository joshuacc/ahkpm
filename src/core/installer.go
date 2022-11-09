package core

import (
	"ahkpm/src/utils"
	"errors"
	"fmt"
	"os"
)

type Installer struct{}

func (i Installer) Install(deps DependencySet) {
	pr := NewPackagesRepository()

	lm, err := LockManifestFromCwd()
	if err == nil && deps.Equals(lm.Dependencies) {
		fmt.Println("No dependency changes found. Installing from lockfile.")
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
	err = resolvedDepTree.ForEach(func(resolvedDepNode TreeNode[ResolvedDependency]) error {
		resolvedDep := resolvedDepNode.Value
		return pr.CopyPackage(resolvedDep, resolvedDep.InstallPath)
	})
	if err != nil {
		utils.Exit(err.Error())
	}

	NewLockManifest().
		WithDependencies(deps).
		WithResolved(resolvedDepTree).
		SaveToCwd()
}

func (i Installer) Update(packageNames ...string) error {
	depsToUpdate := NewDependencySet()
	currentDeps := ManifestFromCwd().Dependencies.AsMap()
	for _, packageName := range packageNames {
		packageName = CanonicalizeDependencyName(packageName)

		dep, ok := currentDeps[packageName]
		if !ok {
			return fmt.Errorf("Cannot update %s. It is not present in ahkpm.json", packageName)
		}
		depsToUpdate.AddDependency(dep)
	}
	if depsToUpdate.Len() != len(packageNames) {
		return errors.New("Cannot update multiple versions of the same package")
	}

	resolver := NewDependencyResolver()
	newResolvedDepTree, err := resolver.Resolve(depsToUpdate)
	if err != nil {
		return err
	}

	// Get tree from lockfile
	lm, err := LockManifestFromCwd()
	if err != nil {
		return err
	}
	oldResolved := ResolvedDependencyTreeFromArray(lm.Resolved)

	// Replace subtrees with new resolved deps
	for _, resolvedDepNode := range newResolvedDepTree {
		for i, oldResolvedDepNode := range oldResolved {
			if resolvedDepNode.Value.Name == oldResolvedDepNode.Value.Name {
				oldResolved[i] = resolvedDepNode
			}
		}
	}

	err = oldResolved.CheckForConflicts()
	if err != nil {
		return err
	}

	// Empty ahkpm-modules
	os.RemoveAll("ahkpm-modules")

	// Copy new resolved deps to ahkpm-modules
	pr := NewPackagesRepository()
	err = oldResolved.ForEach(func(resolvedDepNode TreeNode[ResolvedDependency]) error {
		resolvedDep := resolvedDepNode.Value
		return pr.CopyPackage(resolvedDep, resolvedDep.InstallPath)
	})
	if err != nil {
		return err
	}

	// Save lockfile
	NewLockManifest().
		WithDependencies(lm.Dependencies).
		WithResolved(oldResolved).
		SaveToCwd()

	return nil
}
