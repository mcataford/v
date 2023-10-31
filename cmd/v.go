package main

import (
	"os"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
	subcommands "v/internal/subcommands"
)

// Main entrypoint.
//
// v <command> <...parameters>
//
// v install <version>
//
//	Installs <version> from remote.
//
// v where
//
//	Shows the current Python path.
//
// v use <version>
//
//	Uses the provided version as global python version.
func main() {
	args := os.Args[1:]
	currentState := stateManager.ReadState()

	cli := argparse.CLI{}
	err := cli.AddCommand(
		"where", subcommands.Where, "v where", "Prints the path to the current Python version.",
	).AddCommand(
		"ls", subcommands.ListVersions, "v ls", "Lists the installed Python versions.",
	).AddCommand(
		"install", subcommands.InstallPython, "v install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"which", subcommands.Which, "v which", "Prints the current Python version.",
	).AddCommand(
		"init", subcommands.Initialize, "v init", "Initializes the v state.",
	).AddCommand(
		"use", subcommands.Use, "v use <version>", "Selects which Python version to use.",
	).Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
