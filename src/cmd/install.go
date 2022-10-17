package cmd

import (
	utils "ahkpm/src/utils"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the specified package",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			utils.Exit("Error getting current directory")
		}

		ahkpmFileExists, err := exists(cwd + `\ahkpm.json`)
		if err != nil {
			utils.Exit("Error checking if ahkpm.json exists")
		}

		if !ahkpmFileExists {
			fmt.Println("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
			os.Exit(1)
		}

		if len(args) == 0 {
			// TODO: install all packages in ahkpm.json
			fmt.Println("Please specify a package to install")
			return
		}

		if len(args) > 1 {
			// TODO: support specifying multiple packages
			fmt.Println("Please specify only one package to install")
			return
		}

		packageToInstall := args[0]
		var versionSpecifier string
		if strings.Contains(packageToInstall, "@") {
			splitArg := strings.SplitN(packageToInstall, "@", 2)
			packageToInstall = splitArg[0]
			versionSpecifier = splitArg[1]
		}

		version := getVersion(versionSpecifier)

		installSinglePackage(packageToInstall, version)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

// Installs a single package. For now only support installing an exact version
// TODO: Support omitting version specifier
// TODO: Support version ranges
// TODO: Support branch names
// TODO: Support git commits
func installSinglePackage(packageName string, version Version) {
	fmt.Println("Installing package", packageName, "with", strings.ToLower(string(version.Kind)), version.Value)
	// TODO: validate package name
	// TODO: Halt if missing ahkpm.json

	if version.Kind == SemVerRange || version.Kind == Branch || version.Kind == Commit {
		fmt.Println("Unsupported version type. Ranges, branches, and commits are not yet supported")
		os.Exit(1)
	}

	cacheDir := utils.GetCacheDir()

	packageCacheDir := cacheDir + `\` + packageName

	err := os.MkdirAll(packageCacheDir, os.ModePerm)
	if err != nil {
		utils.Exit("Error creating package cache directory")
	}

	packageWasCloned, err := exists(packageCacheDir + `\.git`)
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
			utils.Exit("Error cloning package")
		}

	} else {
		// Checkout the specified version
		repo, err := git.PlainOpen(packageCacheDir)
		if err != nil {
			utils.Exit("Error opening package")
		}

		worktree, err := repo.Worktree()
		if err != nil {
			utils.Exit("Error getting worktree")
		}

		err = worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewTagReferenceName(version.Value)})
		if err != nil {
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

	// TODO: Add the installed package to ahkpm.json's dependencies list
	// TODO: Create/update a lockfile
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getGitUrl(packageName string) string {
	return "https://" + packageName + ".git"
}

func getVersion(versionSpecifier string) Version {
	v := Version{}

	if utils.IsSemVer(versionSpecifier) {
		v.Kind = SemVerExact
		v.Value = versionSpecifier
	} else if utils.IsSemVerRange(versionSpecifier) {
		v.Kind = SemVerRange
	} else if strings.HasPrefix(versionSpecifier, "branch:") {
		v.Kind = Branch
		v.Value = strings.TrimPrefix(versionSpecifier, "branch:")
	} else if strings.HasPrefix(versionSpecifier, "tag:") {
		v.Kind = Tag
		v.Value = strings.TrimPrefix(versionSpecifier, "tag:")
	} else if strings.HasPrefix(versionSpecifier, "commit:") {
		v.Kind = Commit
		v.Value = strings.TrimPrefix(versionSpecifier, "commit:")
	} else {
		fmt.Println("Invalid version string", versionSpecifier)
		os.Exit(1)
	}

	return v
}

type Version struct {
	Kind  VersionKind
	Value string
}

type VersionKind string

const (
	SemVerExact VersionKind = "SemVerExact"
	SemVerRange VersionKind = "SemVerRange"
	Branch      VersionKind = "Branch"
	Tag         VersionKind = "Tag"
	Commit      VersionKind = "Commit"
)
