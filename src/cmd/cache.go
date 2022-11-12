package cmd

import (
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manipulates the packages cache",
	Long: "Provides subcommands to manipulate the packages cache. The cache is a" +
		" directory where packages are downloaded and stored for later use.",
}

func init() {
	RootCmd.AddCommand(cacheCmd)
}
