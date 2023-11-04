package subcommands

import (
	"os"
	"path"
	argparse "v/argparse"
	stateManager "v/state"
)

var DIRECTORIES = []string{
	"cache",
	"runtimes",
	"shims",
}

var SHIMS = []string{
	"python",
	"python3",
}

const DEFAULT_PERMISSION = 0775

func writeShim(shimPath string) error {
	shimContent := []byte("#!/bin/bash\n$(vm where) $@")
	if err := os.WriteFile(shimPath, shimContent, DEFAULT_PERMISSION); err != nil {
		return err
	}

	return nil
}

// Sets up directories and files used to store downloaded archives,
// installed runtimes and metadata.
func Initialize(args []string, flags argparse.Flags, currentState stateManager.State) error {
	stateDirectory := stateManager.GetStateDirectory()

	os.Mkdir(stateDirectory, DEFAULT_PERMISSION)
	for _, dir := range DIRECTORIES {
		os.Mkdir(stateManager.GetPathFromStateDirectory(dir), DEFAULT_PERMISSION)
	}

	for _, shim := range SHIMS {
		writeShim(stateManager.GetPathFromStateDirectory(path.Join("shims", shim)))
	}

	return nil
}
