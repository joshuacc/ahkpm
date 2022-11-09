package cmd

import (
	"ahkpm/src/core"
	"ahkpm/src/utils"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed install-long.md
var installLong string

var installCmd = &cobra.Command{
	Use:   "install [<packageName>@<version>]...",
	Short: "Installs specified package(s). If none, reinstalls all packages in ahkpm.json.",
	Long:  installLong,
	Run: func(cmd *cobra.Command, args []string) {
		ahkpmFileExists, err := utils.FileExists(`ahkpm.json`)
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
			dependencies := core.ManifestFromCwd().Dependencies
			installer.Install(dependencies)
			return
		}

		deps := make([]core.Dependency, len(args))
		for i, arg := range args {
			dep, err := core.DependencyFromSpecifier(arg)
			if err != nil {
				utils.Exit(err.Error())
			}
			deps[i] = dep
		}

		manifest := core.ManifestFromCwd()
		for _, dep := range deps {
			fmt.Println(
				"Installing package", dep.Name(),
				"with", strings.ToLower(string(dep.Version().Kind())), dep.Version().Value(),
			)
			manifest.Dependencies.AddDependency(dep)
		}

		installer.Install(manifest.Dependencies)
		manifest.SaveToCwd()
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
