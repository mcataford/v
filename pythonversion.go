package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SelectedVersion struct {
	Version string
	Source  string
}

// SearchForPythonVersionFile crawls up to the system root to find any
// .python-version file that could set the current version.
func SearchForPythonVersionFile() (SelectedVersion, bool) {
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

	if versionFound == "" {
		return SelectedVersion{}, false
	}

	return SelectedVersion{Version: versionFound, Source: path.Join(currentPath, ".python-version")}, true
}

// DetermineSelectedPythonVersion returns the Python runtime version that should be
// used according to v.
//
// First, v will look in the current directory and all its parents for a .python-version
// file that would indicate which version is preferred. If none are found, the global
// user-defined version (via `v use <version>`) is used. If there is none, the system
// Python version is used.
func DetermineSelectedPythonVersion(currentState State) (SelectedVersion, error) {
	pythonFileVersion, pythonFileVersionFound := SearchForPythonVersionFile()

	if pythonFileVersionFound {
		return pythonFileVersion, nil
	}

	if len(currentState.GlobalVersion) != 0 {
		return SelectedVersion{Version: currentState.GlobalVersion, Source: GetStatePath("state.json")}, nil
	}

	systemVersion, _ := DetermineSystemPython()
	return SelectedVersion{Source: "system", Version: systemVersion}, nil
}

// DetermineSystemPython returns the unshimmed Python version and path.
// It assumes that /bin/python is where system Python lives.
func DetermineSystemPython() (string, string) {
	versionOut, _ := RunCommand([]string{"/bin/python", "--version"}, GetStatePath(), true)
	detectedVersion, _ := strings.CutPrefix(versionOut, "Python")
	return strings.TrimSpace(detectedVersion), "/bin/python"
}
