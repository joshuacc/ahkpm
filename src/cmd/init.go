package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

type AhkpmJson struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Repository   string            `json:"repository"`
	Website      string            `json:"website"`
	License      string            `json:"license"`
	IssueTracker string            `json:"issueTracker"`
	Author       Person            `json:"author"`
	Dependencies map[string]string `json:"dependencies"`
}

type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

// TODO: Add colors to the prompt
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively create an ahkpm.json file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Unable to get current working directory")
			return
		}
		cwd = filepath.Base(cwd)
		cwd = strings.ToLower(cwd)
		cwd = strings.Replace(cwd, " ", "-", -1)

		// Initialize with default values
		packageSpec := AhkpmJson{
			Name:         cwd,
			Version:      "0.0.1",
			Author:       Person{},
			Dependencies: make(map[string]string),
		}

		var jsonBytes []byte

		for true {
			// TODO: use current directory name as default
			packageSpec.Name = showPrompt(
				"What is the name of your package?",
				prompt.OptionInitialBufferText(packageSpec.Name),
			)

			packageSpec.Version = showPrompt(
				"What version is the package? (Using semantic versioning)",
				prompt.OptionInitialBufferText(packageSpec.Version),
			)

			packageSpec.Description = showPrompt(
				"Please enter a brief description of the package",
				prompt.OptionInitialBufferText(packageSpec.Description),
			)

			// TODO: Look in git config for default
			packageSpec.Repository = showPrompt(
				"What is the URL of the package's git repository? (optional)",
				prompt.OptionInitialBufferText(packageSpec.Repository),
			)

			// TODO: use repo url as default
			packageSpec.Website = showPrompt(
				"What is the URL of the package's homepage? (optional)",
				prompt.OptionInitialBufferText(packageSpec.Website),
			)

			// TODO: use repo url as default
			packageSpec.IssueTracker = showPrompt(
				"What is the URL of the package's bug/issue tracker? (optional)",
				prompt.OptionInitialBufferText(packageSpec.IssueTracker),
			)

			// TODO: use MIT as default
			packageSpec.License = showPrompt(
				"What license is the package released under? (MIT, Apache, etc.) Must be a valid SPDX license identifier.",
				prompt.OptionInitialBufferText(packageSpec.License),
			)

			// TODO: use git config author as default
			packageSpec.Author.Name = showPrompt(
				"What is the author's name/alias? (optional)",
				prompt.OptionInitialBufferText(packageSpec.Author.Name),
			)

			// TODO: use git config email as default
			packageSpec.Author.Email = showPrompt(
				"What is the author's email address? (optional)",
				prompt.OptionInitialBufferText(packageSpec.Author.Email),
			)

			// TODO: use git config author as default
			packageSpec.Author.Website = showPrompt(
				"What is the author's website? (optional)",
				prompt.OptionInitialBufferText(packageSpec.Author.Website),
			)

			jsonBytes, err := json.MarshalIndent(packageSpec, "", "  ")
			if err != nil {
				return
			}

			fmt.Println(string(jsonBytes) + "\n")
			fmt.Println("")

			isCorrect := showPrompt("Is this correct? (y/n)", prompt.OptionInitialBufferText("y"))

			if isCorrect == "y" {
				break
			}

			fmt.Println("Please correct any errors.")
			fmt.Println("")
		}

		ioutil.WriteFile("ahkpm.json", jsonBytes, 0644)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func showPrompt(promptMessage string, options ...prompt.Option) string {
	noSuggestionsCompleter := func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}

	result := ""

	for true {
		fmt.Println(promptMessage)
		result = prompt.Input("> ", noSuggestionsCompleter, options...)
		fmt.Println("")

		if result != "" || strings.Contains(promptMessage, "optional") {
			break
		}
	}

	return result
}

type Validator func(string) (bool string)

func validateNothing(value string) (bool, string) {
	return true, ""
}

func validateRequired(value string) (bool, string) {
	if value == "" {
		return false, "Value is required"
	}

	return true, ""
}
