package cmd

import (
	core "ahkpm/src/core"
	utils "ahkpm/src/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clears the package cache",
	Long:  "Clears the package cache, removing all downloaded packages. Sometimes useful in troubleshooting.",
	Run: func(cmd *cobra.Command, args []string) {
		cleanCache()
	},
}

func init() {
	cacheCmd.AddCommand(cacheCleanCmd)
}

func cleanCache() {
	fmt.Println("Cleaning cache...")
	err := core.NewPackagesRepository().ClearCache()
	if err != nil {
		utils.Exit("Error cleaning cache")
	}

	fmt.Println("Cache cleaned")
}
