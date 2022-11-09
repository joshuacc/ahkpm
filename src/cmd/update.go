package cmd

import (
	core "ahkpm/src/core"
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:        "update <packageName>...",
	SuggestFor: []string{"upgrade"},
	Short:      "Update package(s) to the latest version allowed by ahkpm.json",
	Long: `Updates package(s) to the latest version allowed by ahkpm.json.

For example, if you have a dependency on "github.com/user/repo" with version
"branch:main", running "ahkpm update github.com/user/repo" will update the
package to the latest commit on the main branch.

You may also use package name shorthands, such as "gh:user/repo"`,
	Example: "ahkpm update github.com/joshuacc/mock-ahkpm-package-a",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a package name")
			return
		}
		installer := core.Installer{}
		err := installer.Update(args...)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
