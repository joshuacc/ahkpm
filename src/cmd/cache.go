package cmd

import (
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manipulates the packages cache",
}

func init() {
	RootCmd.AddCommand(cacheCmd)
}
