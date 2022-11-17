package cmd

import (
	"ahkpm/src/core"
	"ahkpm/src/invariant"
	"ahkpm/src/utils"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed install-long.md
var installLong string

//go:embed install-example.txt
var installExample string

var installCmd = &cobra.Command{
	Use:     "install [<packageName>@<version>]...",
	Short:   "Installs specified package(s). If none, reinstalls all packages in ahkpm.json.",
	Long:    installLong,
	Example: installExample,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		ahkpmFileExists, err := utils.FileExists(`ahkpm.json`)
		invariant.AssertNoError(err)

		if !ahkpmFileExists {
			fmt.Println("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
			os.Exit(1)
		}

		installer := core.Installer{}

		newDeps, err := core.NewDependencySet().AddDependenciesFromSpecifiers(args)
		if err != nil {
			utils.Exit(err.Error())
		}

		for _, dep := range newDeps.AsArray() {
			fmt.Println(
				"Installing package", dep.Name(),
				"with", strings.ToLower(string(dep.Version().Kind())), dep.Version().Value(),
			)
		}

		installer.Install(newDeps)
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
