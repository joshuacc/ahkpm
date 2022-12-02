package cmd

import (
	core "ahkpm/src/core"
	utils "ahkpm/src/utils"
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed init-long.md
var initLong string

// TODO: Add colors to the prompt
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively create an ahkpm.json file in the current directory",
	Long:  initLong,
	Run: func(cmd *cobra.Command, args []string) {
		manifestExists, err := utils.FileExists("ahkpm.json")
		if err != nil {
			utils.Exit("Error checking for ahkpm.json")
		}
		if manifestExists {
			utils.Exit("ahkpm.json already exists in this directory")
		}

		initializer := core.Initializer{}
		if cmd.Flag("defaults").Value.String() == "true" {
			initializer.InitFromDefaults()
		} else {
			initializer.InitInteractively()
		}
	},
}

func init() {
	initCmd.Flags().BoolP("defaults", "d", false, "Create an ahkpm.json file with default values. No prompts.")
	RootCmd.AddCommand(initCmd)
}
