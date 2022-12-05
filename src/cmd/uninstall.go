package cmd

import (
	core "ahkpm/src/core"
	utils "ahkpm/src/utils"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <package>...",
	Short:   "Uninstalls the specified package(s)",
	Long:    "Uninstalls the specified package(s)",
	Example: "ahkpm uninstall gh:joshuacc/fake-package",
	Aliases: []string{"remove", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			utils.Exit("Please specify a package name")
		}

		m := core.ManifestFromCwd()

		depNames := make([]string, len(args))
		for i, pkgName := range args {
			pkgName = core.CanonicalizeDependencyName(pkgName)
			if !m.Dependencies.Contains(pkgName) {
				utils.Exit(pkgName + " is not in your dependencies")
			}
			depNames[i] = pkgName
		}

		installer := core.Installer{}
		installer.Uninstall(depNames)
	},
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}
