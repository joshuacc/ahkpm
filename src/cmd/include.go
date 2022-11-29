package cmd

import (
	core "ahkpm/src/core"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed include-long.md
var includeLong string

var includeCmd = &cobra.Command{
	Use:   "include <package>",
	Short: "Gets the \"Include\" statement needed to use a package",
	Long:  includeLong,
	Example: "ahkpm include gh:joshuacc/fake-package\n" +
		"ahkpm include --file my-script.ahk gh:joshuacc/fake-package",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a package name")
			return
		}

		m := core.ManifestFromCwd()
		pkgName := core.CanonicalizeDependencyName(args[0])
		if !m.Dependencies.Contains(pkgName) {
			fmt.Println("Package is not in your dependencies")
			return
		}

		pkgPath := "ahkpm-modules/" + pkgName + "/"
		pkgManifest, err := core.ManifestFromFile(pkgPath + "ahkpm.json")
		if err != nil {
			fmt.Println("Error getting package manifest. It may not have an ahkpm.json file")
			return
		}

		if pkgManifest.Include == "" {
			fmt.Println("Package does not specify which file to \"Include\"")
			return
		}

		includePath := strings.Replace(pkgPath, "/", `\`, -1) + pkgManifest.Include
		includePrefix := "#Include %A_ScriptDir%\\"
		fileName := cmd.Flag("file").Value.String()
		relativePath, err := filepath.Rel(filepath.Dir(fileName), includePath)
		if err != nil {
			fmt.Println("Error getting relative path")
			return
		}
		includeStatement := includePrefix + relativePath
		if fileName == "" {
			fmt.Println(includeStatement)
		} else {
			fileContents, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("Error reading file " + fileName)
				return
			}
			newContents := make([]byte, 0)
			newContents = append(newContents, []byte(includeStatement)...)
			newContents = append(newContents, []byte("\r\n")...)
			newContents = append(newContents, fileContents...)
			err = os.WriteFile(fileName, newContents, 0644)
			if err != nil {
				fmt.Println("Error writing file " + fileName)
				return
			}
		}
	},
}

func init() {
	includeCmd.Flags().StringP("file", "f", "", "The file to add the \"Include\" statement to")
	RootCmd.AddCommand(includeCmd)
}
