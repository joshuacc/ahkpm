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
		core.NewDependency("github.com/joshuacc/mock-ahkpm-package-a", core.NewVersion(core.SemVerRange, "^1.3.1")),
	}

	// Ensure ahkpm.json has the correct dependency and versions
	assert.Equal(t, expectedDeps, m.Dependencies.AsArray())

	lm, err := core.LockManifestFromCwd()
	assert.Nil(t, err)

	expectedResolved := []core.ResolvedDependency{
		{
			Name:        "github.com/joshuacc/mock-ahkpm-package-a",
			Version:     "^1.3.1",
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

	// Ensure that the correct files were installed
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/README.md")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm.json")
	assert.FileExists(t, "ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-a/ahkpm-modules/github.com/joshuacc/mock-ahkpm-package-b/README.md")

	// Cleanup
	err = os.RemoveAll("ahkpm-modules")
	assert.Nil(t, err)
	err = os.Remove("ahkpm.json")
	assert.Nil(t, err)
	err = os.Remove("ahkpm.lock")
	assert.Nil(t, err)
}
