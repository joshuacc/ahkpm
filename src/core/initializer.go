package core

import (
	"ahkpm/src/data"
	"ahkpm/src/invariant"
	"ahkpm/src/utils"
	"fmt"
	"net/mail"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"golang.org/x/exp/slices"
)

type Initializer struct{}

func (i Initializer) InitFromDefaults() {
	GetNewManifestWithDefaults().SaveToCwd()
}

func (i Initializer) InitInteractively() {
	manifest := GetNewManifestWithDefaults()

	for {
		manifest.Version = showPrompt(
			"What version is the package? (Using semantic versioning)",
			validateSemver,
			prompt.OptionInitialBufferText(manifest.Version),
		)

		manifest.Description = showPrompt(
			"Please enter a brief description of the package",
			validateRequired,
			prompt.OptionInitialBufferText(manifest.Description),
		)

		manifest.Repository = showPrompt(
			"What is the URL of the package's git repository? (optional)",
			makeOptional(validateGitHub),
			prompt.OptionInitialBufferText(manifest.Repository),
		)

		// Website defaults to repository url
		if manifest.Website == "" {
			manifest.Website = manifest.Repository
		}

		manifest.Website = showPrompt(
			"What is the URL of the package's homepage? (optional)",
			makeOptional(validateUrl),
			prompt.OptionInitialBufferText(manifest.Website),
		)

		// Issue tracker defaults to GitHub issues if a GitHub repository is specified
		if manifest.IssueTracker == "" && manifest.Repository != "" {
			manifest.IssueTracker = manifest.Repository + "/issues"
		}

		manifest.IssueTracker = showPrompt(
			"What is the URL of the package's bug/issue tracker? (optional)",
			makeOptional(validateUrl),
			prompt.OptionInitialBufferText(manifest.IssueTracker),
		)

		manifest.Include = showPrompt(
			"What is the primary file which users of this package should \"Include\" to use it in their scripts? (optional)",
			validateNothing,
			prompt.OptionInitialBufferText(manifest.Include),
		)

		manifest.License = showPrompt(
			"What license is the package released under? (MIT, Apache, etc.) Must either be a valid SPDX license identifier or \"UNLICENSED\".",
			buildValidatorFromList(data.GetSpdxLicenseIds()),
			prompt.OptionInitialBufferText(manifest.License),
		)

		manifest.Author.Name = showPrompt(
			"What is the author's name/alias? (optional)",
			validateNothing,
			prompt.OptionInitialBufferText(manifest.Author.Name),
		)

		manifest.Author.Email = showPrompt(
			"What is the author's email address? (optional)",
			makeOptional(validateEmail),
			prompt.OptionInitialBufferText(manifest.Author.Email),
		)

		manifest.Author.Website = showPrompt(
			"What is the author's website? (optional)",
			makeOptional(validateUrl),
			prompt.OptionInitialBufferText(manifest.Author.Website),
		)

		fmt.Println(manifest.String() + "\n")
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

	manifest.SaveToCwd()
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
	isMatch, err := regexp.MatchString("^https:\\/\\/github\\.com\\/[\\w-\\.]+\\/[\\w-\\.]+$", value)
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

func GetNewManifestWithDefaults() *Manifest {
	manifest := NewManifest()

	manifest.Version = "1.0.0"
	manifest.Include = getDefaultInclude()
	manifest.Repository = getDefaultRepository()
	manifest.Website = manifest.Repository
	manifest.IssueTracker = getDefaultIssueTracker(manifest.Repository)
	manifest.Include = getDefaultInclude()
	manifest.License = "MIT"
	manifest.Author.Name = getDefaultAuthorName()
	manifest.Author.Email = getDefaultAuthorEmail()

	return manifest
}

func getDefaultRepository() string {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err == nil && string(out) == "true\n" {
		originUrl, err := exec.Command("git", "remote", "get-url", "origin").Output()
		if err == nil && string(originUrl) != "" {
			re, err := regexp.Compile(`(github\.com.+)\.git`)
			if err != nil {
				panic(err)
			}
			submatchResults := re.FindSubmatch(originUrl)
			if len(submatchResults) > 1 {
				domainAndPath := submatchResults[1]
				return "https://" + string(domainAndPath)
			}
		}
	}

	return ""
}

func getDefaultIssueTracker(repository string) string {
	if repository == "" {
		return ""
	}

	return repository + "/issues"
}

// Gets the ahk files in the current working directory. If there is only one,
// it is returned. Otherwise it returns an empty string.
func getDefaultInclude() string {
	ahkFiles, err := filepath.Glob("*.ahk")
	invariant.AssertNoError(err)

	if len(ahkFiles) == 1 {
		return ahkFiles[0]
	}

	return ""
}

func getDefaultAuthorName() string {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err == nil && string(out) == "true\n" {
		userName, err := exec.Command("git", "config", "--get", "user.name").Output()
		if err == nil && string(userName) != "" {
			return strings.Replace(string(userName), "\n", "", -1)
		}
	}

	return ""
}

func getDefaultAuthorEmail() string {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err == nil && string(out) == "true\n" {
		email, err := exec.Command("git", "config", "--get", "user.email").Output()
		if err == nil && string(email) != "" {
			return strings.Replace(string(email), "\n", "", -1)
		}
	}

	return ""
}
