package cmd

import (
	data "ahkpm/src/data"
	utils "ahkpm/src/utils"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
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
			License:      "MIT",
			Author:       Person{},
			Dependencies: make(map[string]string),
		}

		var jsonBytes []byte

		for {
			packageSpec.Name = showPrompt(
				"What is the name of your package?",
				validateNothing,
				prompt.OptionInitialBufferText(packageSpec.Name),
			)

			packageSpec.Version = showPrompt(
				"What version is the package? (Using semantic versioning)",
				validateSemver,
				prompt.OptionInitialBufferText(packageSpec.Version),
			)

			packageSpec.Description = showPrompt(
				"Please enter a brief description of the package",
				validateRequired,
				prompt.OptionInitialBufferText(packageSpec.Description),
			)

			// If we're in a git repository, extract the repository url
			if packageSpec.Repository == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					originUrl, err := exec.Command("git", "remote", "get-url", "origin").Output()
					if err == nil && string(originUrl) != "" {
						re, err := regexp.Compile(`(github\.com.+)\.git`)
						if err != nil {
							panic(err)
						}
						domainAndPath := re.FindSubmatch(originUrl)[1]
						packageSpec.Repository = "https://" + string(domainAndPath)
					}
				}
			}

			packageSpec.Repository = showPrompt(
				"What is the URL of the package's git repository? (optional)",
				makeOptional(validateGitHub),
				prompt.OptionInitialBufferText(packageSpec.Repository),
			)

			// Website defaults to repository url
			if packageSpec.Website == "" {
				packageSpec.Website = packageSpec.Repository
			}

			packageSpec.Website = showPrompt(
				"What is the URL of the package's homepage? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(packageSpec.Website),
			)

			// Issue tracker defaults to GitHub issues if a GitHub repository is specified
			if packageSpec.IssueTracker == "" && packageSpec.Repository != "" {
				packageSpec.IssueTracker = packageSpec.Repository + "/issues"
			}

			packageSpec.IssueTracker = showPrompt(
				"What is the URL of the package's bug/issue tracker? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(packageSpec.IssueTracker),
			)

			packageSpec.License = showPrompt(
				"What license is the package released under? (MIT, Apache, etc.) Must either be a valid SPDX license identifier or \"UNLICENSED\".",
				buildValidatorFromList(data.GetSpdxLicenseIds()),
				prompt.OptionInitialBufferText(packageSpec.License),
			)

			// If we're in a git repository, extract the user's name to use as default
			if packageSpec.Author.Name == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					userName, err := exec.Command("git", "config", "--get", "user.name").Output()
					if err == nil && string(userName) != "" {
						packageSpec.Author.Name = strings.Replace(string(userName), "\n", "", -1)
					}
				}
			}

			packageSpec.Author.Name = showPrompt(
				"What is the author's name/alias? (optional)",
				validateNothing,
				prompt.OptionInitialBufferText(packageSpec.Author.Name),
			)

			// If we're in a git repository, extract the user's name to use as default
			if packageSpec.Author.Email == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					email, err := exec.Command("git", "config", "--get", "user.email").Output()
					if err == nil && string(email) != "" {
						packageSpec.Author.Email = strings.Replace(string(email), "\n", "", -1)
					}
				}
			}

			packageSpec.Author.Email = showPrompt(
				"What is the author's email address? (optional)",
				makeOptional(validateEmail),
				prompt.OptionInitialBufferText(packageSpec.Author.Email),
			)

			packageSpec.Author.Website = showPrompt(
				"What is the author's website? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(packageSpec.Author.Website),
			)

			jsonBytes, err = json.MarshalIndent(packageSpec, "", "  ")
			if err != nil {
				return
			}

			fmt.Println(string(jsonBytes) + "\n")
			fmt.Println("")

			isCorrect := showPrompt(
				"Is this correct? (y/n)",
				validateYesNo,
				prompt.OptionInitialBufferText("y"),
			)

			if isCorrect == "y" {
				break
			}

			fmt.Println("Please correct any errors.")
			fmt.Println("")
		}

		err = os.WriteFile("ahkpm.json", jsonBytes, 0644)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func showPrompt(promptMessage string, validate Validator, options ...prompt.Option) string {
	noSuggestionsCompleter := func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}

	result := ""

	for {
		fmt.Println(promptMessage)
		result = prompt.Input("> ", noSuggestionsCompleter, options...)
		fmt.Println("")

		isValid, message := validate(result)
		if isValid {
			break
		}

		fmt.Println(message)
		fmt.Println("")
	}

	return result
}

type Validator func(string) (isValid bool, message string)

func validateNothing(value string) (bool, string) {
	return true, ""
}

func validateRequired(value string) (bool, string) {
	if value == "" {
		return false, "Value is required"
	}

	return true, ""
}

func validateSemver(value string) (bool, string) {
	isMatch := utils.IsSemVer(value)

	if isMatch {
		return true, ""
	}
	return false, "This is not a valid semantic version. Please see https://semver.org/"
}

func validateGitHub(value string) (bool, string) {
	isMatch, err := regexp.MatchString("^https:\\/\\/github\\.com\\/\\w+\\/\\w+$", value)
	if err != nil {
		panic(err)
	}

	if isMatch {
		return true, ""
	}
	return false, "Please enter a valid GitHub repository URL. Other git hosts will be supported in the future."
}

func validateUrl(value string) (bool, string) {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return false, "This is not a valid URL"
	}
	return true, ""
}

func validateEmail(value string) (bool, string) {
	_, err := mail.ParseAddress(value)
	if err != nil {
		return false, "This is not a valid email address"
	}
	return true, ""
}

func validateYesNo(value string) (bool, string) {
	if value != "y" && value != "n" {
		return false, "Please enter either \"y\" or \"n\""
	}
	return true, ""
}

func buildValidatorFromList(options []string) Validator {
	return func(value string) (bool, string) {
		isInOptions := slices.Contains(options, value)
		if isInOptions {
			return true, ""
		}
		return false, "That is not a valid option"
	}
}

func makeOptional(validator Validator) Validator {
	return func(value string) (bool, string) {
		if value == "" {
			return true, ""
		}

		return validator(value)
	}
}
