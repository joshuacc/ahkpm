package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:        "version",
	Short:      "",
	Deprecated: "Use the --version flag instead. ",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
