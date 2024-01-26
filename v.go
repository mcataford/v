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

	root := cli.Namespace{Label: ""}
	root.AddCommand(
		"init", Initialize, "v init", "Initializes the v state.",
	)

	pythonCommands := cli.Namespace{Label: "python"}
	pythonCommands.AddCommand(
		"install", python.InstallPython, "v python install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"uninstall", python.UninstallPython, "v python uninstall <version>", "Uninstalls the given Python version.",
	).AddCommand(
		"use", python.Use, "v python use <version>", "Selects which Python version to use.",
	).AddCommand(
		"ls", python.ListVersions, "v python ls", "Lists the installed Python versions.",
	).AddCommand(
		"version", python.CurrentVersion, "v python version", "Prints the current version and its source.",
	).AddCommand(
		"which", python.Which, "v python which", "Prints the path to the current Python version.",
	)

	cli := cli.CLI{
		Metadata: map[string]string{
			"Version": Version,
		},
	}

	cli.AddNamespace(root).AddNamespace(pythonCommands)

	err := cli.Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
