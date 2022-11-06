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
	Resolve(deps DependencySet) (resolvedDependencies ResolvedDependencyTree, err error)

	// WithPackagesRepository is used for testing
	WithPackagesRepository(pr PackagesRepository) DependencyResolver
}

type resolver struct {
	packagesRepository PackagesRepository
}

func NewDependencyResolver() DependencyResolver {
	return &resolver{
		packagesRepository: NewPackagesRepository(),
	}
}

func (r *resolver) Resolve(deps DependencySet) (ResolvedDependencyTree, error) {
	if deps.Len() == 0 {
		return ResolvedDependencyTree{}, nil
	}

	resolvedDepNodes, err := r.innerResolve(deps)
	if err != nil {
		return ResolvedDependencyTree{}, err
	}

	depNodesWithInstallPath := resolvedDepNodes.EnsureInstallPaths()

	err = checkForConflicts(depNodesWithInstallPath)
	if err != nil {
		return nil, err
	}

	return depNodesWithInstallPath, nil
}

func (r *resolver) innerResolve(depSet DependencySet) (ResolvedDependencyTree, error) {
	if depSet.Len() == 0 {
		return ResolvedDependencyTree{}, nil
	}

	resolvedDepNodes := make(ResolvedDependencyTree, depSet.Len())

	// For each dependency, get its transitive dependencies.
	for i, dep := range depSet.AsArray() {
		partiallyResolvedDepNode, err := getResolvedDependency(r.packagesRepository, NewTreeNode(dep))
		if err != nil {
			return nil, err
		}

		children, err := r.innerResolve(partiallyResolvedDepNode.Value.Dependencies)
		if err != nil {
			return nil, err
		}

		fullyResolvedDepNode := partiallyResolvedDepNode.WithChildren(children)

		resolvedDepNodes[i] = fullyResolvedDepNode
	}

	return resolvedDepNodes, nil
}

func checkForConflicts(depNodes ResolvedDependencyTree) error {
	allDeps := depNodes.Flatten()

	depMap := make(map[string]ResolvedDependency)
	for _, dep := range allDeps {
		// If the dependency is already in the map, check if the versions are the same.
		if existingDep, ok := depMap[dep.Name]; ok {
			if existingDep.Version != dep.Version {
				return fmt.Errorf("Conflicting versions for dependency %s: %s and %s", dep.Name, existingDep.Version, dep.Version)
			}
			if existingDep.SHA != dep.SHA {
				return fmt.Errorf("Conflicting SHAs for dependency %s: %s and %s", dep.Name, existingDep.SHA, dep.SHA)
			}
		} else {
			depMap[dep.Name] = dep
		}
	}

	return nil
}

func (r *resolver) WithPackagesRepository(pr PackagesRepository) DependencyResolver {
	r.packagesRepository = pr
	return r
}

func getResolvedDependency(pr PackagesRepository, depNode TreeNode[Dependency]) (*TreeNode[ResolvedDependency], error) {
	sha, err := pr.GetResolvedDependencySHA(depNode.Value)
	if err != nil {
		return nil, err
	}

	resolved := ResolvedDependency{
		Name:    depNode.Value.Name(),
		Version: depNode.Value.Version().String(),
		SHA:     sha,
	}

	childDependencies, err := pr.GetPackageDependencies(resolved)
	if err != nil {
		return nil, err
	}

	resolvedNode := NewTreeNode(resolved.WithDependencies(*childDependencies))

	return &resolvedNode, nil
}
