package cmd

import (
	core "ahkpm/src/core"
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <packageName>",
	Short: "Update a package to the latest version which meets the requirements in ahkpm.json",
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
	rootCmd.AddCommand(updateCmd)
}
