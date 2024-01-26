package python

import (
	cli "v/cli"
)

func GetNamespace() cli.Namespace {
	pythonCommands := cli.Namespace{Label: "python"}
	pythonCommands.AddCommand(
		"install", installPython, "v python install <version>", "Downloads, builds and installs a new version of Python.",
	).AddCommand(
		"uninstall", uninstallPython, "v python uninstall <version>", "Uninstalls the given Python version.",
	).AddCommand(
		"use", use, "v python use <version>", "Selects which Python version to use.",
	).AddCommand(
		"ls", listVersions, "v python ls", "Lists the installed Python versions.",
	).AddCommand(
		"version", currentVersion, "v python version", "Prints the current version and its source.",
	).AddCommand(
		"which", which, "v python which", "Prints the path to the current Python version.",
	)

	return pythonCommands
}
