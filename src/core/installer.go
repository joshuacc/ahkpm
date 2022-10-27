package core

import (
	"ahkpm/src/utils"
	"fmt"
	"os"
	"strings"
)

type Installer struct{}

// Installs a single package. For now only support installing an exact version
// TODO: Support omitting version specifier
// TODO: Support version ranges
func (i Installer) InstallSinglePackage(packageName string, version Version) {
	fmt.Println("Installing package", packageName, "with", strings.ToLower(string(version.VersionKind())), version.Value())
	// TODO: validate package name
	hasAhkpmJson, err := utils.FileExists("ahkpm.json")
	if err != nil {
		utils.Exit("Error checking if ahkpm.json exists")
	}
	if !hasAhkpmJson {
		utils.Exit("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
	}

	// Copy files from the package cache to the target module directory
	cwd, err := os.Getwd()
	if err != nil {
		utils.Exit("Error getting current directory")
	}

	pr := NewPackagesRepository()
	dep := NewDependency(packageName, version)
	err = pr.CopyPackage(dep, cwd+`\ahkpm-modules\`+packageName)
	if err != nil {
		utils.Exit(err.Error())
	}

	ManifestFromCwd().
		AddDependency(packageName, version).
		Save()
}

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
