package core

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type ResolvedDependencyTree []TreeNode[ResolvedDependency]

func (r ResolvedDependencyTree) Flatten() []ResolvedDependency {
	allDeps := make([]ResolvedDependency, 0)
	for _, depNode := range r {
		allDeps = append(allDeps, depNode.Flatten()...)
	}
	return allDeps
}

func (r ResolvedDependencyTree) ForEach(callback func(n TreeNode[ResolvedDependency]) error) error {
	for _, depNode := range r {
		err := depNode.ForEach(callback)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r ResolvedDependencyTree) Map(callback func(n TreeNode[ResolvedDependency]) TreeNode[ResolvedDependency]) ResolvedDependencyTree {
	newTree := make(ResolvedDependencyTree, 0)
	for _, depNode := range r {
		newTree = append(newTree, depNode.Map(callback))
	}
	return newTree
}

// EnsureInstallPaths ensures that all resolved dependencies install paths are set correctly
func (r ResolvedDependencyTree) EnsureInstallPaths() ResolvedDependencyTree {
	return r.Map(func(n TreeNode[ResolvedDependency]) TreeNode[ResolvedDependency] {
		n.Value.InstallPath = getRelativeInstallPath(n)
		return n
	})
}

// Merge merges two resolved dependency trees from right to left, by replacing
// the left tree's root nodes with the right tree's root nodes if they have the
// same name. Any root nodes in the right tree that do not exist in the left
// tree are appended to the left tree.
func (r ResolvedDependencyTree) Merge(other ResolvedDependencyTree) ResolvedDependencyTree {
	// convert to map for easier lookup
	rMapToIndex := make(map[string]int, len(r))
	for i, depNode := range r {
		rMapToIndex[depNode.Value.Name] = i
	}

	// merge
	for _, depNode := range other {
		if index, ok := rMapToIndex[depNode.Value.Name]; ok {
			r[index] = depNode
		} else {
			r = append(r, depNode)
		}
	}

	return r
}

func getRelativeInstallPath(n TreeNode[ResolvedDependency]) string {
	path := n.Value.Name
	parent := n.Parent
	for parent != nil {
		path = parent.Value.Name + "/ahkpm-modules/" + path
		parent = parent.Parent
	}

	return "ahkpm-modules/" + path
}

// ResolvedDependencyTreeFromArray takes an array of resolved dependencies (in the format used by LockManifest)
// and converts it into a tree of resolved dependencies
func ResolvedDependencyTreeFromArray(arr []ResolvedDependency) ResolvedDependencyTree {
	// Algorithmically this function is a bit of a mess, but it works and is fast enough

	type intermediateResult struct {
		dep           ResolvedDependency
		dependerNames []string
	}

	// Derive depender names from the install path. The names will be used to build the tree
	tempResults := make([]intermediateResult, len(arr))
	for i, dep := range arr {
		pathWithoutSelf := strings.TrimSuffix(dep.InstallPath, "ahkpm-modules/"+dep.Name)
		pathWithoutPrefix := strings.TrimPrefix(pathWithoutSelf, "ahkpm-modules/")
		pathWithoutEndingSlash := strings.TrimSuffix(pathWithoutPrefix, "/")
		dependerNames := strings.Split(pathWithoutEndingSlash, "/ahkpm-modules/")
		if len(dependerNames) == 1 && dependerNames[0] == "" {
			dependerNames = []string{}
		}
		tempResults[i] = intermediateResult{
			dep:           dep,
			dependerNames: dependerNames,
		}
	}

	tree := make([]TreeNode[ResolvedDependency], 0)

	// Grab the root nodes so that we can attach children to them
	// Iterate over temp results in reverse so that we can remove them as we go
	for i := len(tempResults) - 1; i >= 0; i-- {
		tempResult := tempResults[i]
		if len(tempResult.dependerNames) == 0 {
			tree = append(tree, NewTreeNode(tempResult.dep))
			tempResults = remove(tempResults, i)
		}
	}

	// Iterate over temp results, adding them to the tree and removing them from the queue,
	// stopping when the queue is empty
	for len(tempResults) > 0 {
		// Iterate over temp results in reverse so that we can remove them as we go
		for i := len(tempResults) - 1; i >= 0; i-- {
			tempResult := tempResults[i]
			dependerNode := FindByNamesPath(&tree, tempResult.dependerNames)
			if dependerNode != nil {
				node := NewTreeNode(tempResult.dep)
				dependerNode.AddChild(&node)
				ReplaceByNamesPath(&tree, tempResult.dependerNames, *dependerNode)
				tempResults = remove(tempResults, i)
			}
		}
	}

	SortAllByNames(&tree)

	return tree
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func FindByNamesPath(t *[]TreeNode[ResolvedDependency], names []string) *TreeNode[ResolvedDependency] {
	if len(names) == 0 {
		return nil
	}

	for _, node := range *t {
		if node.Value.Name == names[0] {
			if len(names) == 1 {
				return &node
			}

			return FindByNamesPath(&node.Children, names[1:])
		}
	}

	return nil
}

func ReplaceByNamesPath(t *[]TreeNode[ResolvedDependency], names []string, replacement TreeNode[ResolvedDependency]) {
	if len(names) == 0 {
		return
	}

	for i, node := range *t {
		if node.Value.Name == names[0] {
			if len(names) == 1 {
				(*t)[i] = replacement
				return
			}

			ReplaceByNamesPath(&node.Children, names[1:], replacement)
		}
	}
}

func SortAllByNames(t *[]TreeNode[ResolvedDependency]) {
	slices.SortStableFunc(*t, func(a, b TreeNode[ResolvedDependency]) bool {
		return a.Value.Name < b.Value.Name
	})

	for _, node := range *t {
		SortAllByNames(&node.Children)
	}
}

func (depNodes ResolvedDependencyTree) CheckForConflicts() error {
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
