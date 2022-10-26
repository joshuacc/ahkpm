package core

import (
	"fmt"
)

type DependencyResolver interface {
	// Resolve takes in a list of packages and versions, scans them recursively
	// and returns a tree of all transitive dependencies. If a package occurs
	// more than once, the specified versions are compared. In the case of a conflict,
	// an error is returned. For the time being, any difference in versions is
	// considered a conflict.
	Resolve(deps []Dependency) (resolvedDependencies []TreeNode[Dependency], err error)

	// ReplacePackagesRepository is used for testing
	ReplacePackagesRepository(pr PackagesRepository)
}

type resolver struct {
	packagesRepository PackagesRepository
}

func NewDependencyResolver() DependencyResolver {
	return &resolver{
		packagesRepository: NewPackagesRepository(),
	}
}

func (r *resolver) Resolve(deps []Dependency) ([]TreeNode[Dependency], error) {
	if len(deps) == 0 {
		return []TreeNode[Dependency]{}, nil
	}

	depNodes := make([]TreeNode[Dependency], len(deps))

	// For each dependency, get its transitive dependencies.
	for i, dep := range deps {
		childDependencies, err := r.packagesRepository.GetPackageDependencies(dep)
		if err != nil {
			return nil, err
		}

		children, err := r.Resolve(childDependencies)
		if err != nil {
			return nil, err
		}

		depNodes[i] = TreeNode[Dependency]{
			Value:    dep,
			Children: children,
		}
	}

	err := checkForConflicts(depNodes)
	if err != nil {
		return nil, err
	}

	return depNodes, nil
}

func checkForConflicts(depNodes []TreeNode[Dependency]) error {
	allDeps := make([]Dependency, 0)
	for _, depNode := range depNodes {
		allDeps = append(allDeps, depNode.Flatten()...)
	}

	depMap := make(map[string]Dependency)
	for _, dep := range allDeps {
		// If the dependency is already in the map, check if the versions are the same.
		if existingDep, ok := depMap[dep.Name()]; ok {
			if existingDep.Version() != dep.Version() {
				return fmt.Errorf("conflicting versions for dependency %s: %s and %s", dep.Name(), existingDep.Version(), dep.Version())
			}
		} else {
			depMap[dep.Name()] = dep
		}
	}

	return nil
}

func (r *resolver) ReplacePackagesRepository(pr PackagesRepository) {
	r.packagesRepository = pr
}
