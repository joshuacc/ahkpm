package cmd

import (
	core "ahkpm/src/core"
	"ahkpm/src/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed packages and their versions",
	Long:  `Displays a table listing all installed packages and their versions`,
	Run: func(cmd *cobra.Command, args []string) {
		deps := core.ManifestFromCwd().Dependencies
		fmt.Print(GetDependenciesForDisplay(deps))
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func GetDependenciesForDisplay(set core.DependencySet) string {
	output := ""
	lengthOfLongestName := 0
	for _, dep := range set.AsArray() {
		if len(dep.Name()) > lengthOfLongestName {
			lengthOfLongestName = len(dep.Name())
		}
	}

	lengthOfLongestVersion := 0
	for _, dep := range set.AsArray() {
		if len(dep.Version().String()) > lengthOfLongestVersion {
			lengthOfLongestVersion = len(dep.Version().String())
		}
	}

	output += utils.RightPad("Name", " ", lengthOfLongestName) + "\tVersion\n"
	output += utils.RightPad("", "-", lengthOfLongestName) + "\t"
	output += utils.RightPad("", "-", lengthOfLongestVersion) + "\n"
	for _, dep := range set.AsArray() {
		output = output + utils.RightPad(dep.Name(), " ", lengthOfLongestName) + "\t" + dep.Version().String() + "\n"
	}
	return output
}
