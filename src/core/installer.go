package core

import (
	"ahkpm/src/utils"
	"os"
)

type Installer struct{}

func (i Installer) Install(deps []Dependency) {
	pr := NewPackagesRepository()
	resolver := NewDependencyResolver()
	depTree, err := resolver.Resolve(deps)
	if err != nil {
		utils.Exit(err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		utils.Exit("Error getting current directory")
	}

	for _, dep := range depTree {
		err := dep.ForEach(func(n TreeNode[Dependency]) error {
			path := n.Value.Name()
			parent := n.Parent
			for parent != nil {
				path = parent.Value.Name() + "/ahkpm-modules/" + path
				parent = parent.Parent
			}

			fullPath := cwd + "/ahkpm-modules/" + path

			os.RemoveAll(fullPath)
			err := pr.CopyPackage(n.Value, fullPath)
			if err != nil {
				return err
			}
			err = pr.CopyPackage(n.Value, fullPath)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			utils.Exit(err.Error())
		}
	}
}

func getGitUrl(packageName string) string {
	return "https://" + packageName + ".git"
}
