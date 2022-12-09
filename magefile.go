//go:build mage
// +build mage

package main

import (
	"ahkpm/src/cmd"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/princjef/mageutil/bintool"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		"-coverpkg", "./...",
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
	out, err := sh.Output("bin/ahkpm.exe", "--version")
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

func MarkdownDocs(dir string) error {
	return generateAhkpmMarkdownTree(cmd.RootCmd, dir)
}

func generateAhkpmMarkdownTree(cmd *cobra.Command, dir string) error {
	md := generateMarkdownDocs(cmd)

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"
	filename := filepath.Join(dir, basename)
	err := os.WriteFile(filename, []byte(md), 0644)
	if err != nil {
		return err
	}

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		err := generateAhkpmMarkdownTree(c, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateMarkdownDocs(cmd *cobra.Command) string {
	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	yamlFrontMatterTemplate := `---
title: %s
description: %s
lead: %s
menu:
  docs:
    parent: "commands"
toc: true
---
`
	buf.WriteString(fmt.Sprintf(yamlFrontMatterTemplate, name, cmd.Short, cmd.Short))

	if (cmd.Long != "") && (cmd.Long != cmd.Short) {
		buf.WriteString("## Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	buf.WriteString("## Usage\n\n")
	buf.WriteString("```text\n")
	if cmd.HasAvailableSubCommands() {
		buf.WriteString(cmd.CommandPath() + " [command]\n")
	}
	buf.WriteString(cmd.UseLine() + "\n")
	buf.WriteString("```\n\n")

	if len(cmd.Example) > 0 {
		buf.WriteString("## Examples\n\n")
		buf.WriteString(fmt.Sprintf("```text\n%s\n```\n\n", cmd.Example))
	}

	if cmd.HasAvailableSubCommands() {
		buf.WriteString("## Available subcommands\n\n")
		for _, c := range cmd.Commands() {
			buf.WriteString(fmt.Sprintf("- `%s`: %s\n", c.Name(), c.Short))
		}
		buf.WriteString("\n")
	}

	if cmd.HasAvailableFlags() {
		buf.WriteString("## Options\n\n")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			buf.WriteString(fmt.Sprintf("- `--%s`, `-%s`: %s\n", f.Name, f.Shorthand, f.Usage))
		})
	}

	if len(cmd.Aliases) > 0 {
		buf.WriteString("## Aliases\n\n")
		buf.WriteString(fmt.Sprintf("`%s`\n\n", strings.Join(cmd.Aliases, "`, `")))
	}

	return buf.String()
}
