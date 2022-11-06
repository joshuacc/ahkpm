package cmd

import (
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

		// installer := core.Installer{}
		// installer.Update(args...)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
