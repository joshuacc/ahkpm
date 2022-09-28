//go:build mage
// +build mage

package main

import (
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

	return linter.Command("run").Run()
}

func Test() error {
	return sh.RunV("go", "test", "-timeout", "30s", "-cover", "./...")
}

func Build() error {
	return sh.Run("go", "build", "-o", "bin/ahkpm.exe", "./src")
}

func CIVerify() error {
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
