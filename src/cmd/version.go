package cmd

import (
	"ahkpm/src/constants"
	utils "ahkpm/src/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of ahkpm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("     ahkpm: " + constants.SelfVersion)
		version, err := utils.GetAutoHotkeyVersion()
		if err == nil {
			fmt.Println("AutoHotkey: " + version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
