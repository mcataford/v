package main

import (
	"io/ioutil"
	"os"
	"path"
	"slices"
	"strings"
)

// SearchForPythonVersionFile crawls up to the system root to find any
// .python-version file that could set the current version.
func SearchForPythonVersionFile() (string, bool) {
	currentPath, _ := os.Getwd()
	var versionFound string
	for {
		content, err := ioutil.ReadFile(path.Join(currentPath, ".python-version"))

		if err == nil {
			versionFound = strings.TrimSpace(string(content))
			break
		}

		nextPath := path.Dir(currentPath)

		if currentPath == nextPath {
			break
		}

		currentPath = nextPath
	}

	return versionFound, versionFound != ""
}

// DetermineSelectedPythonVersion returns the Python runtime version that should be
// used according to v.
//
// First, v will look in the current directory and all its parents for a .python-version
// file that would indicate which version is preferred. If none are found, the global
// user-defined version (via `v use <version>`) is used. If there is none, the system
// Python version is used.
func DetermineSelectedPythonVersion(currentState State) (string, error) {
	pythonFileVersion, pythonFileVersionFound := SearchForPythonVersionFile()

	if pythonFileVersionFound {
		return pythonFileVersion, nil
	}

	if len(currentState.GlobalVersion) != 0 {
		return currentState.GlobalVersion, nil
	}

	return "SYSTEM", nil
}

// DetermineSystemPython returns the unshimmed Python version and path.
// This is done by inspected the output of `which` and `python --version` if v's shims
// are not in $PATH.
func DetermineSystemPython() (string, string) {
	currentPathEnv := os.Getenv("PATH")
	pathWithoutShims := slices.DeleteFunc(strings.Split(currentPathEnv, ":"), func(element string) bool {
		return element == GetStatePath("shims")
	})

	// FIXME: This should be set through RunCommand instead.
	os.Setenv("PATH", strings.Join(pathWithoutShims, ":"))
	defer os.Setenv("PATH", currentPathEnv)

	whichOut, _ := RunCommand([]string{"which", "python"}, GetStatePath(), true)
	versionOut, _ := RunCommand([]string{"python", "--version"}, GetStatePath(), true)

	detectedVersion, _ := strings.CutPrefix(versionOut, "Python")

	return strings.TrimSpace(detectedVersion), strings.TrimSpace(whichOut)
}
