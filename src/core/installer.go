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

	NewManifest().
		ReadFromCwd().
		AddDependency(packageName, version).
		Save()
	// TODO: Create/update a lockfile
}

func getGitUrl(packageName string) string {
	return "https://" + packageName + ".git"
}
