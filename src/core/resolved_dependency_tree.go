package core

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

func (r ResolvedDependencyTree) EnsureInstallPaths() ResolvedDependencyTree {
	return r.Map(func(n TreeNode[ResolvedDependency]) TreeNode[ResolvedDependency] {
		n.Value.InstallPath = getRelativeInstallPath(n)
		return n
	})
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
