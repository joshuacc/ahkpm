package core

import (
	"ahkpm/src/utils"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/otiai10/copy"
)

type Installer struct{}

// Installs a single package. For now only support installing an exact version
// TODO: Support omitting version specifier
// TODO: Support version ranges
func (i Installer) InstallSinglePackage(packageName string, version Version) {
	fmt.Println("Installing package", packageName, "with", strings.ToLower(string(version.Kind)), version.Value)
	// TODO: validate package name
	hasAhkpmJson, err := utils.FileExists("ahkpm.json")
	if err != nil {
		utils.Exit("Error checking if ahkpm.json exists")
	}
	if !hasAhkpmJson {
		utils.Exit("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
	}

	cacheDir := utils.GetCacheDir()

	packageCacheDir := cacheDir + `\` + packageName

	err = os.MkdirAll(packageCacheDir, os.ModePerm)
	if err != nil {
		utils.Exit("Error creating package cache directory")
	}

	packageWasCloned, err := utils.FileExists(packageCacheDir + `\.git`)
	if err != nil {
		utils.Exit("Error checking if package was cloned")
	}

	if !packageWasCloned {
		// Clone the repository into the cache directory
		_, err := git.PlainClone(packageCacheDir, false, &git.CloneOptions{
			URL:               getGitUrl(packageName),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			message := "Error cloning package"
			if err.Error() == "authentication required" {
				message = "Error downloading package " + packageName + ". Are you sure that package exists?"
			}
			utils.Exit(message)
		}
	}

	// Checkout the specified version
	repo, err := git.PlainOpen(packageCacheDir)
	if err != nil {
		utils.Exit("Error opening package")
	}

	worktree, err := repo.Worktree()
	if err != nil {
		utils.Exit("Error getting worktree")
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(version.Value))
	if err != nil {
		message := "Error resolving revision"
		if err.Error() == "reference not found" {
			message = "Could not find version " + version.String() + " for package " + packageName + ". Are you sure that version exists?"
		}
		utils.Exit(message)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:  (*hash),
		Force: true, // Ignore changes in the working tree
	})
	if err != nil {
		fmt.Println(err.Error())
		utils.Exit("Error checking out version")
	}

	submodules, err := worktree.Submodules()
	if err != nil {
		utils.Exit("Error getting submodules")
	}

	for _, sub := range submodules {
		err := sub.Update(&git.SubmoduleUpdateOptions{})
		if err != nil {
			utils.Exit("Error updating submodule")
		}
	}

	// Copy files from the package cache to the target module directory
	cwd, err := os.Getwd()
	if err != nil {
		utils.Exit("Error getting current directory")
	}

	targetModuleDir := cwd + `\ahkpm-modules\` + packageName

	err = copy.Copy(packageCacheDir, targetModuleDir)
	if err != nil {
		utils.Exit("Error copying package to target module directory")
	}

	AhkpmJson{}.
		ReadFromFile().
		AddDependency(packageName, version).
		Save()
	// TODO: Create/update a lockfile
}

func getGitUrl(packageName string) string {
	return "https://" + packageName + ".git"
}
