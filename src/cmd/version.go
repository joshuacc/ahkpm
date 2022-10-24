package cmd

import (
	"ahkpm/src/constants"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of ahkpm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ahkpm version: " + constants.SelfVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
