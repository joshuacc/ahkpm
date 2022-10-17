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
		Exit("Error validating semver:")
	}

	return isMatch
}

// TODO: Implement this
func IsSemVerRange(value string) bool {
	return false
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func GetCacheDir() string {
	value, succeeded := os.LookupEnv("userprofile")
	if !succeeded {
		fmt.Println("Unable to get userprofile")
		os.Exit(1)
	}
	return value + `\.ahkpm`
}
