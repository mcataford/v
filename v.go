package main

import (
	"os"
	cli "v/cli"
	python "v/python"
	state "v/state"
)

const (
	Version  = "0.0.8"
	Author   = "Marc Cataford <hello@karnov.club>"
	Homepage = "https://github.com/mcataford/v"
)

// Main entrypoint.
func main() {
	args := os.Args[1:]
	currentState := state.ReadState()

	cli := cli.CLI{
		Metadata: map[string]string{
			"Version": Version,
		},
	}

	err := cli.AddCommand(
		"install", python.InstallPython, "v install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"uninstall", python.UninstallPython, "v uninstall <version>", "Uninstalls the given Python version.",
	).AddCommand(
		"use", python.Use, "v use <version>", "Selects which Python version to use.",
	).AddCommand(
		"ls", python.ListVersions, "v ls", "Lists the installed Python versions.",
	).AddCommand(
		"version", python.CurrentVersion, "v version", "Prints the current version and its source.",
	).AddCommand(
		"which", python.Which, "v which", "Prints the path to the current Python version.",
	).AddCommand(
		"init", Initialize, "v init", "Initializes the v state.",
	).Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
