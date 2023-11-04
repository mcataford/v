package main

import (
	"os"
	argparse "v/argparse"
	stateManager "v/state"
	subcommands "v/commands"
)

const (
	Version = "0.0.1"
	Author  = "mcataford <hello@karnov.club>"
)

// Main entrypoint.
func main() {
	args := os.Args[1:]
	currentState := stateManager.ReadState()

	cli := argparse.CLI{
		Metadata: map[string]string{
			"Version": Version,
		},
	}

	err := cli.AddCommand(
		"install", subcommands.InstallPython, "v install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"uninstall", subcommands.UninstallPython, "v uninstall <version>", "Uninstalls the given Python version.",
	).AddCommand(
		"use", subcommands.Use, "v use <version>", "Selects which Python version to use.",
	).AddCommand(
		"ls", subcommands.ListVersions, "v ls", "Lists the installed Python versions.",
	).AddCommand(
		"where", subcommands.Where, "v where", "Prints the path to the current Python version.",
	).AddCommand(
		"which", subcommands.Which, "v which", "Prints the current Python version.",
	).AddCommand(
		"init", subcommands.Initialize, "v init", "Initializes the v state.",
	).Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
