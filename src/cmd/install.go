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
		if len(args) == 0 {
			// TODO: install all packages
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
// TODO: Support no version specified
// TODO: Support version ranges
// TODO: Support branch names
// TODO: Support git tags
// TODO: Support git commits
func installSinglePackage(packageName string, version Version) {
	fmt.Println("Installing package", packageName, "with", strings.ToLower(string(version.Kind)), version.Value)
	// TODO: validate package name
	// TODO: Halt if missing ahkpm.json

	if version.Kind == SemVerRange || version.Kind == Branch || version.Kind == Commit {
		fmt.Println("Unsupported version type. Ranges, branches, and commits are not yet supported")
		os.Exit(1)
	}

	cacheDir := getCacheDir()

	packageCacheDir := cacheDir + `\` + packageName

	os.MkdirAll(packageCacheDir, os.ModePerm)

	packageWasCloned, err := exists(packageCacheDir + `\.git`)
	if err != nil {
		fmt.Println("Error checking if package was cloned", err)
		os.Exit(1)
	}

	if !packageWasCloned {
		// Clone the repository into the cache directory
		_, err := git.PlainClone(packageCacheDir, false, &git.CloneOptions{URL: getGitUrl(packageName)})
		if err != nil {
			fmt.Println("Error cloning package", err)
			os.Exit(1)
		}
	} else {
		// Checkout the specified version
		repo, err := git.PlainOpen(packageCacheDir)
		if err != nil {
			fmt.Println("Error opening package", err)
			os.Exit(1)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println("Error getting worktree", err)
			os.Exit(1)
		}

		err = worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewTagReferenceName(version.Value)})
		if err != nil {
			fmt.Println("Error checking out version", err)
			os.Exit(1)
		}
	}

	// Copy files from the package cache to the target module directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory", err)
		os.Exit(1)
	}

	targetModuleDir := cwd + `\ahkpm-modules\` + packageName

	copy.Copy(packageCacheDir, targetModuleDir)

	// TODO: Add the installed package to ahkpm.json's dependencies list
	// TODO: Create/update a lockfile
}

func getCacheDir() string {
	value, succeeded := os.LookupEnv("userprofile")
	if !succeeded {
		fmt.Println("Unable to get userprofile")
		os.Exit(1)
	}
	return value + `\.ahkpm`
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
