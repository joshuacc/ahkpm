package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs the `test` script from ahkpm.json if it exists",
	Long: "Runs the `test` script from ahkpm.json if it exists. This is a" +
		" convenience command equivalent to `ahkpm run test`.",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		RunScript("test")
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
