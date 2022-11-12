package cmd

import (
	core "ahkpm/src/core"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed update-long.md
var updateLong string

var updateCmd = &cobra.Command{
	Use:        "update <packageName>...",
	SuggestFor: []string{"upgrade"},
	Short:      "Update package(s) to the latest version allowed by ahkpm.json",
	Long:       updateLong,
	Example:    "ahkpm update github.com/joshuacc/fake-package\nahkpm update gh:joshuacc/fake-package",
	Run: func(cmd *cobra.Command, args []string) {
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

func init() {
	RootCmd.AddCommand(updateCmd)
}
