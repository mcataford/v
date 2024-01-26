package main

import (
	"os"
	"slices"
	cli "v/cli"
	logger "v/logger"
	python "v/python"
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
	shimContent := []byte("#!/bin/bash\n$(v which --raw) $@")
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
func UninstallPython(args []string, flags cli.Flags, currentState state.State) error {
	runtimePath := state.GetStatePath("runtimes", "py-"+args[1])
	err := os.RemoveAll(runtimePath)
	return err
}

func InstallPython(args []string, flags cli.Flags, currentState state.State) error {
	version := args[1]

	return python.InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
}

func Use(args []string, flags cli.Flags, currentState state.State) error {
	version := args[1]
	if err := python.ValidateVersion(version); err != nil {
		return err
	}

	availableVersions := state.GetAvailableVersions()

	found := false
	for _, v := range availableVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		logger.InfoLogger.Println("Version not installed. Installing it first.")
		python.InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
	}

	state.WriteState(version)
	logger.InfoLogger.Printf("Now using Python %s\n", version)

	return nil
}
func ListVersions(args []string, flags cli.Flags, currentState state.State) error {
	installedVersions, err := python.ListInstalledVersions()

	if err != nil {
		return err
	}

	if len(installedVersions) == 0 {
		logger.InfoLogger.Println("No versions installed!")
		return nil
	}

	for _, d := range installedVersions {
		logger.InfoLogger.Println(d)
	}

	return nil
}

// Which prints out the system path to the executable being used by `python`.
func Which(args []string, flags cli.Flags, currentState state.State) error {
	selectedVersion, _ := python.DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := python.ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	var printedPath string

	if selectedVersion.Source == "system" {
		_, sysPath := python.DetermineSystemPython()
		printedPath = sysPath + " (system)"
	} else if isInstalled {
		tag := python.VersionStringToStruct(selectedVersion.Version)
		printedPath = state.GetStatePath("runtimes", "py-"+selectedVersion.Version, "bin", "python"+tag.MajorMinor())
	} else {
		logger.InfoLogger.Printf("The desired version (%s) is not installed.\n", selectedVersion.Version)
		return nil
	}

	prefix := "Python path: "

	if flags.RawOutput {
		prefix = ""
	} else {
		printedPath = logger.Bold(printedPath)
	}

	logger.InfoLogger.Printf("%s%s\n", prefix, printedPath)
	return nil
}

// CurrentVersion (called via `v version`) outputs the currently selected version
// and what configures it. If the version is configured by a file, the file is returned
// under "source", if the system Python is used, "system" is returned as a source.
func CurrentVersion(args []string, flags cli.Flags, currentState state.State) error {
	selectedVersion, _ := python.DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := python.ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	if !isInstalled {
		logger.InfoLogger.Println(logger.Bold(logger.Yellow("WARNING: This version is not installed.")))
	}

	if flags.RawOutput {
		logger.InfoLogger.Println(selectedVersion.Version)
		return nil
	}

	logger.InfoLogger.Printf("Python version: %s\nSource: %s\n", logger.Bold(selectedVersion.Version), logger.Bold(selectedVersion.Source))
	return nil
}
