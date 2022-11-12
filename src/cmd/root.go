package cmd

import (
	"ahkpm/src/constants"
	utils "ahkpm/src/utils"
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed root-long.md
var rootLong string

var RootCmd = &cobra.Command{
	Use:   "ahkpm",
	Short: "The root command for the ahkpm CLI",
	Long:  rootLong,
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
	DisableAutoGenTag: true,
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
