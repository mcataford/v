package main

import (
	"fmt"
	"os"
	"slices"
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
func Initialize(args []string, flags Flags, currentState State) error {
	if flags.AddPath {
		fmt.Printf("export PATH=%s:$PATH\n", GetStatePath("shims"))
		return nil
	}

	os.Mkdir(GetStatePath(), DEFAULT_PERMISSION)
	for _, dir := range DIRECTORIES {
		os.Mkdir(GetStatePath(dir), DEFAULT_PERMISSION)
	}

	for _, shim := range SHIMS {
		writeShim(GetStatePath("shims", shim))
	}

	return nil
}
func UninstallPython(args []string, flags Flags, currentState State) error {
	runtimePath := GetStatePath("runtimes", fmt.Sprintf("py-%s", args[1]))
	err := os.RemoveAll(runtimePath)
	return err
}

func InstallPython(args []string, flags Flags, currentState State) error {
	version := args[1]

	return InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
}

func Use(args []string, flags Flags, currentState State) error {
	version := args[1]
	if err := ValidateVersion(version); err != nil {
		return err
	}

	availableVersions := GetAvailableVersions()

	found := false
	for _, v := range availableVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Version not installed. Installing it first.")
		InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
	}

	WriteState(version)
	fmt.Printf("Now using Python %s\n", version)

	return nil
}
func ListVersions(args []string, flags Flags, currentState State) error {
	installedVersions, err := ListInstalledVersions()

	if err != nil {
		return err
	}

	if len(installedVersions) == 0 {
		fmt.Println("No versions installed!")
		return nil
	}

	for _, d := range installedVersions {
		fmt.Println(d)
	}

	return nil
}

// Which prints out the system path to the executable being used by `python`.
func Which(args []string, flags Flags, currentState State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	var printedPath string

	if selectedVersion.Source == "system" {
		_, sysPath := DetermineSystemPython()
		printedPath = fmt.Sprintf("%s (system)", sysPath)
	} else if isInstalled {
		tag := VersionStringToStruct(selectedVersion.Version)
		printedPath = GetStatePath("runtimes", fmt.Sprintf("py-%s", selectedVersion.Version), "bin", fmt.Sprintf("python%s", tag.MajorMinor()))
	} else {
		fmt.Printf("The desired version (%s) is not installed.\n", selectedVersion.Version)
		return nil
	}

	prefix := "Python path: "

	if flags.RawOutput {
		prefix = ""
	} else {
		printedPath = Bold(printedPath)
	}

	fmt.Printf("%s%s\n", prefix, printedPath)
	return nil
}

// CurrentVersion (called via `v version`) outputs the currently selected version
// and what configures it. If the version is configured by a file, the file is returned
// under "source", if the system Python is used, "system" is returned as a source.
func CurrentVersion(args []string, flags Flags, currentState State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	if !isInstalled {
		fmt.Println(Bold(Yellow("WARNING: This version is not installed.")))
	}

	if flags.RawOutput {
		fmt.Println(selectedVersion.Version)
		return nil
	}

	fmt.Printf("Python version: %s\nSource: %s\n", Bold(selectedVersion.Version), Bold(selectedVersion.Source))
	return nil
}
