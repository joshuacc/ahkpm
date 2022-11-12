package cmd

import (
	core "ahkpm/src/core"
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

	output += rightPad("Name", " ", lengthOfLongestName) + "\tVersion\n"
	output += rightPad("", "-", lengthOfLongestName) + "\t"
	output += rightPad("", "-", lengthOfLongestVersion) + "\n"
	for _, dep := range set.AsArray() {
		output = output + rightPad(dep.Name(), " ", lengthOfLongestName) + "\t" + dep.Version().String() + "\n"
	}
	return output
}

func rightPad(s string, char string, length int) string {
	for len(s) < length {
		s += char
	}
	return s
}
