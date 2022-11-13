package cmd

import (
	core "ahkpm/src/core"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed update-long.md
var updateLong string
var all bool

var updateCmd = &cobra.Command{
	Use:        "update <packageName>...",
	SuggestFor: []string{"upgrade"},
	Short:      "Update package(s) to the latest version allowed by ahkpm.json",
	Long:       updateLong,
	Example:    "ahkpm update github.com/joshuacc/fake-package\nahkpm update gh:joshuacc/fake-package",
	Run: func(cmd *cobra.Command, args []string) {
		deps := core.ManifestFromCwd().Dependencies
		packages := GetDependencies(deps)
		if all {
			installer := core.Installer{}
			err := installer.Update(packages...)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		if len(args) == 0 {
			fmt.Println("Please specify a package name")
			return
		}
		installer := core.Installer{}
		err := installer.Update(args...)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func GetDependencies(set core.DependencySet) []string {
	var allPackages []string
	for _, dep := range set.AsArray() {
		allPackages = append(allPackages, dep.Name())
	}
	return allPackages
}

func init() {
	updateCmd.Flags().BoolVar(&all, "all", false, "updates all dependencies")
	RootCmd.AddCommand(updateCmd)
}
