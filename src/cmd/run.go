package cmd

import (
	"ahkpm/src/core"
	utils "ahkpm/src/utils"
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//go:embed run-long.md
var runLong string

var runCmd = &cobra.Command{
	Use:     "run <script>",
	Short:   "Run user-defined scripts from ahkpm.json",
	Long:    runLong,
	Example: `ahkpm run build\n` + `ahkpm run test`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			utils.Exit("No script name provided")
		}
		RunScript(args[0])
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func RunScript(scriptName string) {
	script, ok := core.ManifestFromCwd().Scripts[scriptName]
	if !ok {
		utils.Exit(fmt.Sprintf("Script '%s' not found in ahkpm.json", scriptName))
	}

	fmt.Println("> " + script)

	scriptCmd := exec.Command("pwsh", "-c", script)

	// Allow the script to use the current program's stdin, stdout, and stderr
	scriptCmd.Stdout = os.Stdout
	scriptCmd.Stderr = os.Stderr
	scriptCmd.Stdin = os.Stdin

	err := scriptCmd.Run()
	if err != nil {
		utils.Exit(fmt.Sprintf("Script \"%s\" failed with %s", scriptName, err.Error()))
	}
}
