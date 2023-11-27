package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
const DEFAULT_SYSTEM_PY_PATH = "/bin/python"

func writeShim(shimPath string) error {
	shimContent := []byte("#!/bin/bash\n$(v where) $@")
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
		return errors.New("Version not installed.")
	}

	WriteState(version)
	fmt.Printf("Now using Python %s\n", version)

	return nil
}
func ListVersions(args []string, flags Flags, currentState State) error {
	runtimesDir := GetStatePath("runtimes")
	entries, err := os.ReadDir(runtimesDir)

	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No versions installed!")
		return nil
	}

	for _, d := range entries {
		fmt.Println(strings.TrimPrefix(d.Name(), "py-"))
	}

	return nil
}

// Where prints out the system path to the executable being used by `python`.
func Where(args []string, flags Flags, currentState State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)

	if selectedVersion == "SYSTEM" {
		_, sysPath := DetermineSystemPython()
		fmt.Println(sysPath)
		return nil
	}

	tag := VersionStringToStruct(selectedVersion)
	fmt.Println(GetStatePath("runtimes", fmt.Sprintf("py-%s", selectedVersion), "bin", fmt.Sprintf("python%s", tag.MajorMinor())))
	return nil
}

// Which prints out the Python version that will be used by shims. It can be invoked
// directly by the `v which` command.
//
// If no version is set (i.e. none is installed, the specified version is not installed),
// the system version is used and 'SYSTEM' is printed by Which.
func Which(args []string, flags Flags, currentState State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)

	if selectedVersion == "SYSTEM" {
		sysVersion, _ := DetermineSystemPython()
		fmt.Println(sysVersion)
		return nil
	}

	fmt.Println(selectedVersion)
	return nil
}
