package core

import (
	"ahkpm/src/utils"
	"errors"
	"fmt"
	"os"
)

type Installer struct{}

func (i Installer) Install(newDeps DependencySet) {
	pr := NewPackagesRepository()

	lm, err := LockManifestFromCwd()
	if err == nil && newDeps.Len() == 0 {
		fmt.Println("No dependency changes found. Installing from lockfile.")
		os.RemoveAll("ahkpm-modules")
		for _, resolvedDep := range lm.Resolved {
			err := pr.CopyPackage(resolvedDep, resolvedDep.InstallPath)
			if err != nil {
				utils.Exit(err.Error())
			}
		}

		fmt.Println("Installation complete.")
		return
	}
	hasLockfile := err == nil

	manifest := ManifestFromCwd()

	// If there is no lockfile, we need to resolve all dependencies, not just
	// the new ones.
	deps := newDeps
	if !hasLockfile {
		for _, dep := range manifest.Dependencies.AsArray() {
			deps.AddDependency(dep)
		}
	}

	resolver := NewDependencyResolver()
	resolvedDepTree, err := resolver.Resolve(deps)
	if err != nil {
		utils.Exit(err.Error())
	}

	var combinedDepTree ResolvedDependencyTree
	if hasLockfile {
		oldDepTree := ResolvedDependencyTreeFromArray(lm.Resolved)
		combinedDepTree = oldDepTree.Merge(resolvedDepTree)
	} else {
		combinedDepTree = resolvedDepTree
	}

	err = combinedDepTree.CheckForConflicts()
	if err != nil {
		utils.Exit(err.Error())
	}

	i.copyResolved(combinedDepTree)

	manifest.Dependencies.AddDependencies(newDeps.AsArray())
	manifest.SaveToCwd()

	NewLockManifest().
		WithDependencies(manifest.Dependencies).
		WithResolved(combinedDepTree).
		SaveToCwd()

	fmt.Println("Installation complete.")
}

func (i Installer) Uninstall(depNames []string) {
	manifest := ManifestFromCwd()
	manifest.Dependencies.RemoveDependenciesByName(depNames)

	lm, err := LockManifestFromCwd()
	if err != nil {
		utils.Exit(err.Error())
	}

	resolvedDepTree := ResolvedDependencyTreeFromArray(lm.Resolved).RemoveTopLevelDependencies(depNames)

	i.copyResolved(resolvedDepTree)

	manifest.SaveToCwd()

	NewLockManifest().
		WithDependencies(manifest.Dependencies).
		WithResolved(resolvedDepTree).
		SaveToCwd()

	fmt.Println("Uninstallation complete.")
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

	i.copyResolved(oldResolved)

	// Save lockfile
	NewLockManifest().
		WithDependencies(lm.Dependencies).
		WithResolved(oldResolved).
		SaveToCwd()

	return nil
}

func (i Installer) copyResolved(resolved ResolvedDependencyTree) {
	os.RemoveAll("ahkpm-modules")
	err := resolved.ForEach(func(resolvedDepNode TreeNode[ResolvedDependency]) error {
		resolvedDep := resolvedDepNode.Value
		return NewPackagesRepository().CopyPackage(resolvedDep, resolvedDep.InstallPath)
	})
	if err != nil {
		utils.Exit(err.Error())
	}
}
