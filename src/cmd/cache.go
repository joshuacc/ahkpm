package cmd

import (
	core "ahkpm/src/core"
	utils "ahkpm/src/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache clean",
	Short: "Clears the package cache",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a subcommand")
			return
		}

		if args[0] == "clean" {
			cleanCache()
		}
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}

func cleanCache() {
	fmt.Println("Cleaning cache...")
	err := core.NewPackagesRepository().ClearCache()
	if err != nil {
		utils.Exit("Error cleaning cache")
	}

	fmt.Println("Cache cleaned")
}
