package cmd

import (
	core "ahkpm/src/core"
	"ahkpm/src/invariant"
	"ahkpm/src/utils"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

//go:embed version-long.md
var versionLong string

var versionCmd = &cobra.Command{
	Use:     "version [<newversion> | major | minor | patch]",
	Short:   "Bumps the version in `ahkpm.json`.",
	Long:    versionLong,
	Example: "ahkpm version major",
	Run: func(cmd *cobra.Command, args []string) {
		ver := args[0]
		manifest := core.ManifestFromCwd()
		oldVersion := manifest.Version
		newVersion, err := GetUpdatedVersion(oldVersion, ver)
		if err != nil {
			utils.Exit(err.Error())
		}

		manifest.Version = newVersion
		manifest.SaveToCwd()

		fmt.Println("Version bumped from " + oldVersion + " to " + newVersion)
	},
}

func init() {
	versionCmd.Flags().StringP("message", "m", "", "Custom git commit and tag message")
	RootCmd.AddCommand(versionCmd)
}

func GetUpdatedVersion(currentVersion string, newVersionSpec string) (string, error) {
	current, err := semver.StrictNewVersion(currentVersion)
	invariant.AssertNoError(err)

	if newVersionSpec == "major" {
		return current.IncMajor().String(), nil
	} else if newVersionSpec == "minor" {
		return current.IncMinor().String(), nil
	} else if newVersionSpec == "patch" {
		return current.IncPatch().String(), nil
	} else if utils.IsSemVer(newVersionSpec) {
		return newVersionSpec, nil
	}

	return "", errors.New("Invalid version")
}
