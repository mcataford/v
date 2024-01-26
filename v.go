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

	cli := cli.CLI{
		Metadata: map[string]string{
			"Version": Version,
		},
	}

	cli.AddNamespace(root).AddNamespace(python.GetNamespace())

	err := cli.Run(args, currentState)

	if err != nil {
		panic(err)
	}
}
