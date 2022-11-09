package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:        "version",
	Short:      "",
	Hidden:     true,
	Deprecated: "Use the --version flag instead. ",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
