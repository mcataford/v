package main

import (
	"os"
	cli "v/cli"
	logger "v/logger"
	state "v/state"
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
	shimContent := []byte("#!/bin/bash\n$(v python which --raw) $@")
	if err := os.WriteFile(shimPath, shimContent, DEFAULT_PERMISSION); err != nil {
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

	os.Mkdir(state.GetStatePath(), DEFAULT_PERMISSION)
	for _, dir := range DIRECTORIES {
		os.Mkdir(state.GetStatePath(dir), DEFAULT_PERMISSION)
	}

	for _, shim := range SHIMS {
		writeShim(state.GetStatePath("shims", shim))
	}

	return nil
}
