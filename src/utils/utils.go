package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func StructFromFile[T any](path string, s *T) (*T, error) {
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error reading " + path)
	}
	err = json.Unmarshal(jsonBytes, s)
	if err != nil {
		panic(err)
		//return nil, errors.New("Error unmarshalling " + path)
	}
	return s, nil
}

func GetAhkpmDir() string {
	value, succeeded := os.LookupEnv("userprofile")
	if !succeeded {
		fmt.Println("Unable to get userprofile")
		os.Exit(1)
	}
	return value + `\.ahkpm`
}

func GetAutoHotkeyVersion() (string, error) {
	versionScript := `FileAppend, %A_AhkVersion%, *`
	scriptPath := GetAhkpmDir() + `\version.ahk`
	err := os.WriteFile(scriptPath, []byte(versionScript), 0644)
	if err != nil {
		return "", err
	}

	// The pipe to more is needed to get the output to stdout
	out, err := exec.Command("autohotkey", scriptPath, "| more").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
