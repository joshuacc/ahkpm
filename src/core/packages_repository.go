package core

import (
	"ahkpm/src/utils"
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/otiai10/copy"
)

type PackagesRepository interface {
	CopyPackage(dep Dependency, path string) error
	GetPackageDependencies(dep Dependency) ([]Dependency, error)
	ClearCache() error
}

type packagesRepository struct{}

func NewPackagesRepository() PackagesRepository {
	return &packagesRepository{}
}

func (pr *packagesRepository) CopyPackage(dep Dependency, path string) error {
	err := pr.ensurePackageIsReady(dep)
	if err != nil {
		return err
	}
	err = copy.Copy(pr.getPackageCacheDir(dep), path)
	if err != nil {
		return errors.New("Error copying package to target module directory")
	}
	return nil
}

func (pr *packagesRepository) GetPackageDependencies(dep Dependency) ([]Dependency, error) {
	return []Dependency{}, nil
}

func (pr *packagesRepository) ClearCache() error {
	return os.RemoveAll(pr.getCacheDir())
}

func (pr *packagesRepository) getCacheDir() string {
	value, succeeded := os.LookupEnv("userprofile")
	if !succeeded {
		utils.Exit("Unable to get userprofile")
	}
	return value + `\.ahkpm\cache`
}

func (pr *packagesRepository) getPackageCacheDir(dep Dependency) string {
	return pr.getCacheDir() + `\` + dep.Name()
}

func (pr *packagesRepository) ensurePackageIsReady(dep Dependency) error {
	packageCacheDir := pr.getPackageCacheDir(dep)

	err := os.MkdirAll(packageCacheDir, os.ModePerm)
	if err != nil {
		return errors.New("Error creating package cache directory")
	}

	packageWasCloned, err := utils.FileExists(packageCacheDir + `\.git`)
	if err != nil {
		return errors.New("Error checking if package was cloned")
	}

	if !packageWasCloned {
		// Clone the repository into the cache directory
		_, err := git.PlainClone(packageCacheDir, false, &git.CloneOptions{
			URL:               getGitUrl(dep.Name()),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			message := "Error cloning package"
			if err.Error() == "authentication required" {
				message = "Error downloading package " + dep.Name() + ". Are you sure that package exists?"
			}
			return errors.New(message)
		}
	}

	// Checkout the specified version
	repo, err := git.PlainOpen(packageCacheDir)
	if err != nil {
		return errors.New("Error opening package")
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return errors.New("Error getting worktree")
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(dep.Version().Value()))
	if err != nil {
		message := "Error resolving revision"
		if err.Error() == "reference not found" {
			message = "Could not find version " + dep.Version().String() + " for package " + dep.Name() + ". Are you sure that version exists?"
		}
		return errors.New(message)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:  (*hash),
		Force: true, // Ignore changes in the working tree
	})
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Error checking out version")
	}

	submodules, err := worktree.Submodules()
	if err != nil {
		return errors.New("Error getting submodules")
	}

	for _, sub := range submodules {
		err := sub.Update(&git.SubmoduleUpdateOptions{})
		if err != nil {
			return errors.New("Error updating submodule")
		}
	}
	return nil
}
