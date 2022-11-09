package cmd

import (
	"ahkpm/src/constants"
	utils "ahkpm/src/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ahkpm",
	Short: "The package manager for AutoHotkey",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			fmt.Print(GetVersions())
			return
		}
		err := cmd.Help()
		if err != nil {
			utils.Exit(err.Error())
		}
	},
}

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Display the version of ahkpm and AutoHotkey")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetVersions() string {
	versions := "     ahkpm: " + constants.SelfVersion + "\n"
	ahkVersion, err := utils.GetAutoHotkeyVersion()
	if err == nil {
		versions = versions + "AutoHotkey: " + ahkVersion + "\n"
	}
	return versions
}
