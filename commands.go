package main

import (
	"maps"
	"os"
	cli "v/cli"
	logger "v/logger"
	python "v/python"
	state "v/state"
)

var stateDirectories = []string{
	"cache",
	"runtimes",
	"shims",
}

const defaultFilePermissions = 0775

func writeShim(shimPath string, shimCall string) error {
	shimContent := []byte("#!/bin/bash\n" + shimCall)
	if err := os.WriteFile(shimPath, shimContent, defaultFilePermissions); err != nil {
		return err
	}

	return nil
}

// Sets up directories and files used to store downloaded archives,
// installed runtimes and metadata.
func Initialize(args []string, flags cli.Flags, currentState state.State) error {
	if flags.AddPath {
		logger.InfoLogger.Printf("export PATH=%s:$PATH\n", state.GetStatePath("shims"))
		return nil
	}

	os.Mkdir(state.GetStatePath(), defaultFilePermissions)
	logger.InfoLogger.Printf("Created state directory: %s\n", state.GetStatePath())
	for _, dir := range stateDirectories {
		newPath := state.GetStatePath(dir)
		os.Mkdir(newPath, defaultFilePermissions)
		logger.InfoLogger.Printf("Created %s\n", newPath)
	}

	allShims := map[string]string{}

	maps.Copy(allShims, python.Shims)

	for shimName, shimContent := range allShims {
		newShim := state.GetStatePath("shims", shimName)
		writeShim(newShim, shimContent)
		logger.InfoLogger.Printf("Created shim: %s\n", newShim)
	}

	return nil
}
