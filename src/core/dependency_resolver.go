package core

import "ahkpm/src/utils"

type DependencyResolver interface {
	// Resolve takes in a list of packages and versions, scans them recursively
	// and returns a tree of all transitive dependencies. If a package occurs
	// more than once, the specified versions are compared. In the case of a conflict,
	// an error is returned. For the time being, any difference in versions is
	// considered a conflict.
	Resolve(deps []Dependency) (resolvedDependencies []Node[Dependency], err error)

	// ReplacePackagesRepository is used for testing
	ReplacePackagesRepository(pr PackagesRepository)
}

type Node[T any] struct {
	Value    T
	Children []Node[T]
}

type resolver struct {
	packagesRepository PackagesRepository
}

func NewDependencyResolver() DependencyResolver {
	return &resolver{
		packagesRepository: NewPackagesRepository(),
	}
}

func (r *resolver) Resolve(deps []Dependency) ([]Node[Dependency], error) {
	if len(deps) == 0 {
		return []Node[Dependency]{}, nil
	}

	depNodes := make([]Node[Dependency], len(deps))

	// For each dependency, get its transitive dependencies.
	for i, dep := range deps {
		childDependencies, err := r.packagesRepository.GetPackageDependencies(dep)
		if err != nil {
			utils.Exit(err.Error())
		}

		children, err := r.Resolve(childDependencies)
		if err != nil {
			return nil, err
		}

		depNodes[i] = Node[Dependency]{
			Value:    dep,
			Children: children,
		}
	}

	return depNodes, nil
}

func (r *resolver) ReplacePackagesRepository(pr PackagesRepository) {
	r.packagesRepository = pr
}

func ArrayToNodes[T any](arr []T) []Node[T] {
	nodes := make([]Node[T], len(arr))
	for i, item := range arr {
		nodes[i] = Node[T]{Value: item}
	}
	return nodes
}
