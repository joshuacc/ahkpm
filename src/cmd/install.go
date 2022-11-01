package cmd

import (
	"ahkpm/src/core"
	"ahkpm/src/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs either the specified package or all packages listed in ahkpm.json",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			utils.Exit("Error getting current directory")
		}

		ahkpmFileExists, err := utils.FileExists(cwd + `\ahkpm.json`)
		if err != nil {
			utils.Exit("Error checking if ahkpm.json exists")
		}

		if !ahkpmFileExists {
			fmt.Println("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
			os.Exit(1)
		}

		installer := core.Installer{}

		if len(args) == 0 {
			fmt.Println("Installing all dependencies")
			dependencies := core.ManifestFromCwd().Dependencies()
			installer.Install(dependencies)
			return
		}

		if len(args) > 1 {
			// TODO: support specifying multiple packages
			fmt.Println("Please specify only one package to install")
			return
		}

		newDep, err := core.DependencyFromSpecifier(args[0])
		if err != nil {
			utils.Exit(err.Error())
		}

		fmt.Println(
			"Installing package", newDep.Name(),
			"with", strings.ToLower(string(newDep.Version().Kind())), newDep.Version().Value(),
		)
		manifest := core.ManifestFromCwd()
		manifest.AddDependency(newDep)
		installer.Install(manifest.Dependencies())
		manifest.SaveToCwd()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
