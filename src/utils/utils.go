package utils

import (
	"fmt"
	"os"
	"regexp"
)

func IsSemVer(value string) bool {
	// This regular expression is taken from semver.org
	isMatch, err := regexp.MatchString("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$", value)
	if err != nil {
		fmt.Println("Error validating semver:", err)
		os.Exit(1)
	}

	return isMatch
}

// TODO: Implement this
func IsSemVerRange(value string) bool {
	return false
}
