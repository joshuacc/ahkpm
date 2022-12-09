package main_test

import (
	. "ahkpm/src/cmd"
	"ahkpm/src/core"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallingNewPackageWithDependency(t *testing.T) {
	core.NewManifest().SaveToCwd()
	RootCmd.SetArgs([]string{"install", "gh:joshuacc/mock-ahkpm-package-a"})
	err := RootCmd.Execute()
	assert.Nil(t, err)

	m := core.ManifestFromCwd()

	expectedDeps := []core.Dependency{
		core.NewDependency("github.com/joshuacc/mock-ahkpm-package-a", core.NewVersion(core.SemVerRange, "^1.3.2")),
	}

	// Ensure ahkpm.json has the correct dependency and versions
	assert.Equal(t, expectedDeps, m.Dependencies.AsArray())

	lm, err := core.LockManifestFromCwd()
	assert.Nil(t, err)

	expectedResolved := []core.ResolvedDependency{
		{
			Name:        "github.com/joshuacc/mock-ahkpm-package-a",
			Version:     "^1.3.2",
			SHA:         "a7ecc280bf13fd81a4b28fef373b4e2b311f265f",
			InstallPath: "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a",
			Dependencies: core.NewDependencySet().
				AddDependency(
					core.NewDependency("github.com/joshuacc/mock-ahkpm-package-b", core.NewVersion(core.Branch, "main")),
				),
		},

		{
			Name:         "github.com/joshuacc/mock-ahkpm-package-b",
			Version:      "branch:main",
			SHA:          "c4ada0b84f91a7e673fc4bd687e805154adb67d5",
			InstallPath:  "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b",
			Dependencies: core.NewDependencySet(),
		},
	}

	// Ensure that ahkpm.lock has the correct list of resolved dependencies, including transitive dependencies
	assert.Equal(t, expectedResolved, lm.Resolved)

	// Ensure that the correct files were installed
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/README.md")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm.json")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b/README.md")

	cleanupFiles(t)
}

func TestUninstallingPackageWithDependency(t *testing.T) {
	core.NewManifest().SaveToCwd()
	RootCmd.SetArgs([]string{"install", "gh:joshuacc/mock-ahkpm-package-a"})
	err := RootCmd.Execute()
	assert.Nil(t, err)

	// Ensure that the correct files were installed
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/README.md")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm.json")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b/README.md")

	RootCmd.SetArgs([]string{"uninstall", "gh:joshuacc/mock-ahkpm-package-a"})
	err = RootCmd.Execute()
	assert.Nil(t, err)

	m := core.ManifestFromCwd()

	// Ensure ahkpm.json has no dependencies listed
	assert.Equal(t, 0, len(m.Dependencies.AsArray()))

	lm, err := core.LockManifestFromCwd()
	assert.Nil(t, err)

	// Ensure that ahkpm.lock has no resolved dependencies
	assert.Equal(t, 0, len(lm.Resolved))

	// Ensure that the correct files were uninstalled
	assert.NoFileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/README.md")
	assert.NoFileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm.json")
	assert.NoFileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b/README.md")

	cleanupFiles(t)
}

func TestInstallingFromLockFile(t *testing.T) {
	core.NewManifest().SaveToCwd()
	RootCmd.SetArgs([]string{"install", "gh:joshuacc/mock-ahkpm-package-a"})
	err := RootCmd.Execute()
	assert.Nil(t, err)

	err = os.RemoveAll("ahkpm-modules")
	assert.Nil(t, err)
	assert.NoFileExists(t, "ahkpm-modules")

	RootCmd.SetArgs([]string{"install"})
	err = RootCmd.Execute()
	assert.Nil(t, err)

	// Ensure that the correct files were installed
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/README.md")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm.json")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b/README.md")

	cleanupFiles(t)
}

func TestUpdating(t *testing.T) {
	core.NewManifest().SaveToCwd()
	RootCmd.SetArgs([]string{"install", "gh:joshuacc/mock-ahkpm-package-a@1.3.1"})
	err := RootCmd.Execute()
	assert.Nil(t, err)

	m := core.ManifestFromCwd()

	expectedDeps := []core.Dependency{
		core.NewDependency("github.com/joshuacc/mock-ahkpm-package-a", core.NewVersion(core.SemVerExact, "1.3.1")),
	}

	// Ensure ahkpm.json has the correct dependency and versions
	assert.Equal(t, expectedDeps, m.Dependencies.AsArray())

	lm, err := core.LockManifestFromCwd()
	assert.Nil(t, err)

	expectedResolved := []core.ResolvedDependency{
		{
			Name:        "github.com/joshuacc/mock-ahkpm-package-a",
			Version:     "1.3.1",
			SHA:         "9d750cfbcffa05b8a33590a4a66c8047fd057452",
			InstallPath: "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a",
			Dependencies: core.NewDependencySet().
				AddDependency(
					core.NewDependency("github.com/joshuacc/mock-ahkpm-package-b", core.NewVersion(core.Branch, "main")),
				),
		},
		{
			Name:         "github.com/joshuacc/mock-ahkpm-package-b",
			Version:      "branch:main",
			SHA:          "c4ada0b84f91a7e673fc4bd687e805154adb67d5",
			InstallPath:  "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b",
			Dependencies: core.NewDependencySet(),
		},
	}

	// Ensure that ahkpm.lock has the correct list of resolved dependencies, including transitive dependencies
	assert.Equal(t, expectedResolved, lm.Resolved)

	err = os.RemoveAll("ahkpm-modules")
	assert.Nil(t, err)

	newDepWithRange := core.NewDependency("github.com/joshuacc/mock-ahkpm-package-a", core.NewVersion(core.SemVerRange, "^1.3.1"))
	m.Dependencies.AddDependency(newDepWithRange)
	m.SaveToCwd()

	lm.Dependencies.AddDependency(newDepWithRange)
	lm.SaveToCwd()

	RootCmd.SetArgs([]string{"update"})
	err = UpdateCmd.Flags().Set("all", "true")
	assert.Nil(t, err)
	err = RootCmd.Execute()
	assert.Nil(t, err)

	lm, err = core.LockManifestFromCwd()
	assert.Nil(t, err)

	expectedResolved = []core.ResolvedDependency{
		{
			Name:        "github.com/joshuacc/mock-ahkpm-package-a",
			Version:     "^1.3.1",
			SHA:         "a7ecc280bf13fd81a4b28fef373b4e2b311f265f", // 1.3.2
			InstallPath: "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a",
			Dependencies: core.NewDependencySet().
				AddDependency(
					core.NewDependency("github.com/joshuacc/mock-ahkpm-package-b", core.NewVersion(core.Branch, "main")),
				),
		},
		{
			Name:         "github.com/joshuacc/mock-ahkpm-package-b",
			Version:      "branch:main",
			SHA:          "c4ada0b84f91a7e673fc4bd687e805154adb67d5",
			InstallPath:  "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b",
			Dependencies: core.NewDependencySet(),
		},
	}

	// Ensure that ahkpm.lock has the correct list of resolved dependencies, including transitive dependencies
	assert.Equal(t, expectedResolved, lm.Resolved)

	cleanupFiles(t)
}

func TestInitDefaults(t *testing.T) {
	RootCmd.SetArgs([]string{"init"})
	err := InitCmd.Flags().Set("defaults", "true")
	assert.Nil(t, err)
	err = RootCmd.Execute()
	assert.Nil(t, err)

	m := core.ManifestFromCwd()

	// Ensure that the manifest has the correct default values
	assert.Equal(t, "1.0.0", m.Version)
	assert.Equal(t, "", m.Include)
	assert.Equal(t, "MIT", m.License)

	err = os.RemoveAll("ahkpm.json")
	assert.Nil(t, err)
}

func TestRunScripts(t *testing.T) {
	// Ensure file doesn't exist before running the script
	assert.NoFileExists(t, "test.txt")

	m := core.NewManifest()
	m.Scripts["writefile"] = "echo test > test.txt"
	m.SaveToCwd()
	RootCmd.SetArgs([]string{"run", "writefile"})
	err := RootCmd.Execute()
	assert.Nil(t, err)

	assert.FileExists(t, "test.txt")

	err = os.Remove("ahkpm.json")
	assert.Nil(t, err)
	err = os.Remove("test.txt")
	assert.Nil(t, err)
}

func cleanupFiles(t *testing.T) {
	err := os.RemoveAll("ahkpm-modules")
	assert.Nil(t, err)
	err = os.Remove("ahkpm.json")
	assert.Nil(t, err)
	err = os.Remove("ahkpm.lock")
	assert.Nil(t, err)
}
