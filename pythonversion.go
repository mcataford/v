package main

import (
	"os"
	"slices"
	"strings"
)

// DetermineSelectedPythonVersion returns the Python runtime version that should be
// used according to v.
//
// By default, 'SYSTEM' is returned, which signals that the non-v-managed Python
// runtime is used.
func DetermineSelectedPythonVersion(currentState State) (string, error) {
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
	whichOut, _ := RunCommand([]string{"which", "python"}, ".", true)
	versionOut, _ := RunCommand([]string{"python", "--version"}, ".", true)

	detectedVersion, _ := strings.CutPrefix(versionOut, "Python ")

	return detectedVersion, whichOut
}
