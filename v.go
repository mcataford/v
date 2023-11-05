package main

import (
	"os"
)

const (
	Version = "0.0.1"
	Author  = "mcataford <hello@karnov.club>"
)

// Main entrypoint.
func main() {
	args := os.Args[1:]
	currentState := ReadState()

	cli := CLI{
		Metadata: map[string]string{
			"Version": Version,
		},
	}

	err := cli.AddCommand(
		"install", InstallPython, "v install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"uninstall", UninstallPython, "v uninstall <version>", "Uninstalls the given Python version.",
	).AddCommand(
		"use", Use, "v use <version>", "Selects which Python version to use.",
	).AddCommand(
		"ls", ListVersions, "v ls", "Lists the installed Python versions.",
	).AddCommand(
		"where", Where, "v where", "Prints the path to the current Python version.",
	).AddCommand(
		"which", Which, "v which", "Prints the current Python version.",
	).AddCommand(
		"init", Initialize, "v init", "Initializes the v state.",
	).Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
