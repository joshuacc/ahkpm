package cmd

import (
	core "ahkpm/src/core"
	data "ahkpm/src/data"
	utils "ahkpm/src/utils"
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

// TODO: Add colors to the prompt
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively create an ahkpm.json file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Abort init if there is already an ahkpm.json file in the current directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Unable to get current working directory")
			return
		}
		cwd = filepath.Base(cwd)
		cwd = strings.ToLower(cwd)
		cwd = strings.Replace(cwd, " ", "-", -1)

		// Initialize with default values
		aj := core.AhkpmJson{
			Name:    cwd,
			Version: "0.0.1",
			License: "MIT",
		}.New()

		for {
			aj.Name = showPrompt(
				"What is the name of your package?",
				validateNothing,
				prompt.OptionInitialBufferText(aj.Name),
			)

			aj.Version = showPrompt(
				"What version is the package? (Using semantic versioning)",
				validateSemver,
				prompt.OptionInitialBufferText(aj.Version),
			)

			aj.Description = showPrompt(
				"Please enter a brief description of the package",
				validateRequired,
				prompt.OptionInitialBufferText(aj.Description),
			)

			// If we're in a git repository, extract the repository url
			if aj.Repository == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					originUrl, err := exec.Command("git", "remote", "get-url", "origin").Output()
					if err == nil && string(originUrl) != "" {
						re, err := regexp.Compile(`(github\.com.+)\.git`)
						if err != nil {
							panic(err)
						}
						domainAndPath := re.FindSubmatch(originUrl)[1]
						aj.Repository = "https://" + string(domainAndPath)
					}
				}
			}

			aj.Repository = showPrompt(
				"What is the URL of the package's git repository? (optional)",
				makeOptional(validateGitHub),
				prompt.OptionInitialBufferText(aj.Repository),
			)

			// Website defaults to repository url
			if aj.Website == "" {
				aj.Website = aj.Repository
			}

			aj.Website = showPrompt(
				"What is the URL of the package's homepage? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(aj.Website),
			)

			// Issue tracker defaults to GitHub issues if a GitHub repository is specified
			if aj.IssueTracker == "" && aj.Repository != "" {
				aj.IssueTracker = aj.Repository + "/issues"
			}

			aj.IssueTracker = showPrompt(
				"What is the URL of the package's bug/issue tracker? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(aj.IssueTracker),
			)

			aj.License = showPrompt(
				"What license is the package released under? (MIT, Apache, etc.) Must either be a valid SPDX license identifier or \"UNLICENSED\".",
				buildValidatorFromList(data.GetSpdxLicenseIds()),
				prompt.OptionInitialBufferText(aj.License),
			)

			// If we're in a git repository, extract the user's name to use as default
			if aj.Author.Name == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					userName, err := exec.Command("git", "config", "--get", "user.name").Output()
					if err == nil && string(userName) != "" {
						aj.Author.Name = strings.Replace(string(userName), "\n", "", -1)
					}
				}
			}

			aj.Author.Name = showPrompt(
				"What is the author's name/alias? (optional)",
				validateNothing,
				prompt.OptionInitialBufferText(aj.Author.Name),
			)

			// If we're in a git repository, extract the user's name to use as default
			if aj.Author.Email == "" {
				out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
				if err == nil && string(out) == "true\n" {
					email, err := exec.Command("git", "config", "--get", "user.email").Output()
					if err == nil && string(email) != "" {
						aj.Author.Email = strings.Replace(string(email), "\n", "", -1)
					}
				}
			}

			aj.Author.Email = showPrompt(
				"What is the author's email address? (optional)",
				makeOptional(validateEmail),
				prompt.OptionInitialBufferText(aj.Author.Email),
			)

			aj.Author.Website = showPrompt(
				"What is the author's website? (optional)",
				makeOptional(validateUrl),
				prompt.OptionInitialBufferText(aj.Author.Website),
			)

			fmt.Println(aj.String() + "\n")
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

		aj.Save()
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
