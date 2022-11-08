//go:build mage
// +build mage

package main

import (
	"errors"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/princjef/mageutil/bintool"
)

var linter = bintool.Must(bintool.New(
	"golangci-lint{{.BinExt}}",
	"1.49.0",
	"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
))

func Lint() error {
	if err := linter.Ensure(); err != nil {
		return err
	}

	return linter.Command("run --timeout=5m").Run()
}

func Test() error {
	return sh.RunV(
		"go", "test",
		"-covermode", "atomic",
		"-coverprofile", "coverage.out",
		"-timeout", "30s",
		"-cover",
		"./...",
	)
}

func Build() error {
	return sh.Run("go", "build", "-o", "bin/ahkpm.exe", "./src")
}

func BuildWithVersion(version string) error {
	err := os.WriteFile("src/constants/ahkpm-version.txt", []byte(version), 0644)
	if err != nil {
		return err
	}
	return Build()
}

func Verify() error {
	if err := Lint(); err != nil {
		return err
	}

	if err := Test(); err != nil {
		return err
	}

	if err := Build(); err != nil {
		return err
	}

	return nil
}

// This command requires both Wix and go-msi to be installed on the system.
// GitHub actions has Wix installed by default, but go-msi must be installed.
func Msi(version string) error {
	return sh.Run("go-msi", "make", "--src", "wix-templates", "--msi", "bin/ahkpm-"+version+".msi", "--version", version)
}

func VerifyVersion(version string) error {
	out, err := sh.Output("bin/ahkpm.exe", "version")
	if err != nil {
		return err
	}

	if !strings.Contains(out, version) {
		return errors.New("Version mismatch.\n    Expected: " + version + "\n    Actual: " + out)
	}

	return nil
}

func BuildRelease(version string) error {
	if err := BuildWithVersion(version); err != nil {
		return err
	}

	if err := VerifyVersion(version); err != nil {
		return err
	}

	if err := Msi(version); err != nil {
		return err
	}
	return nil
}
