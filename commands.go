package main

import (
	"errors"
	"fmt"
	"os"
	"path"
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

func writeShim(shimPath string) error {
	shimContent := []byte("#!/bin/bash\n$(vm where) $@")
	if err := os.WriteFile(shimPath, shimContent, DEFAULT_PERMISSION); err != nil {
		return err
	}

	return nil
}

// Sets up directories and files used to store downloaded archives,
// installed runtimes and metadata.
func Initialize(args []string, flags Flags, currentState State) error {
	stateDirectory := GetStateDirectory()

	os.Mkdir(stateDirectory, DEFAULT_PERMISSION)
	for _, dir := range DIRECTORIES {
		os.Mkdir(GetPathFromStateDirectory(dir), DEFAULT_PERMISSION)
	}

	for _, shim := range SHIMS {
		writeShim(GetPathFromStateDirectory(path.Join("shims", shim)))
	}

	return nil
}
func UninstallPython(args []string, flags Flags, currentState State) error {
	runtimePath := GetPathFromStateDirectory(path.Join("runtimes", fmt.Sprintf("py-%s", args[1])))
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
	runtimesDir := GetPathFromStateDirectory("runtimes")
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

func Where(args []string, flags Flags, currentState State) error {
	version := currentState.GlobalVersion
	tag := VersionStringToStruct(version)
	withoutPatch := fmt.Sprintf("%s.%s", tag.Major, tag.Minor)
	fmt.Printf("%s/runtimes/py-%s/bin/python%s\n", GetStateDirectory(), currentState.GlobalVersion, withoutPatch)
	return nil
}

func Which(args []string, flags Flags, currentState State) error {
	fmt.Println(currentState.GlobalVersion)
	return nil
}