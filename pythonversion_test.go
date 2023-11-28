package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// SetupAndCleanupEnvironment sets up a test directory and
// environment variables before the test and returns a cleanup
// function that can be deferred to cleanup any changes to the
// system.
func SetupAndCleanupEnvironment(t *testing.T) func() {
	os.Setenv("V_ROOT", t.TempDir())

	temporaryWd := t.TempDir()

	os.Chdir(temporaryWd)

	return func() {
		os.Unsetenv("V_ROOT")
	}
}

func TestDetermineSystemPythonGetsUnshimmedPythonRuntime(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	ioutil.WriteFile(GetStatePath("shims", "python"), []byte("#!/bin/bash\necho \"Python 4.5.6\""), 0777)
	mockSystemPythonPath := t.TempDir()
	mockSystemPythonExecPath := path.Join(mockSystemPythonPath, "python")
	ioutil.WriteFile(mockSystemPythonExecPath, []byte("#!/bin/bash\necho \"Python 1.2.3\""), 0777)

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s:/usr/bin", GetStatePath("shims"), mockSystemPythonPath))
	defer os.Setenv("PATH", oldPath)
	sysVersion, sysPath := DetermineSystemPython()

	if sysVersion != "1.2.3" {
		t.Errorf("Expected system Python to be 1.2.3, found %s instead.", sysVersion)
	}

	if sysPath != mockSystemPythonExecPath {
		t.Errorf("Expected system Python path to be %s, found %s instead.", mockSystemPythonExecPath, sysPath)
	}
}

func TestDetermineSelectedPythonVersionUsesPythonVersionFileIfFound(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	// Writing a mock user-defined state.
	mockState := State{GlobalVersion: "1.0.0"}
	statePath := GetStatePath("state.json")
	stateData, _ := json.Marshal(mockState)
	ioutil.WriteFile(statePath, stateData, 0750)

	temporaryWd := t.TempDir()
	os.Chdir(temporaryWd)
	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)

	version, err := DetermineSelectedPythonVersion(ReadState())

	if err != nil || version != "1.2.3" {
		t.Errorf("Expected version to be %s, got %s instead.", "1.2.3", version)
	}
}

func TestDetermineSelectedPythonVersionGetsUserDefinedVersion(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	// Writing a mock user-defined state.
	mockState := State{GlobalVersion: "1.0.0"}
	statePath := GetStatePath("state.json")
	stateData, _ := json.Marshal(mockState)
	ioutil.WriteFile(statePath, stateData, 0750)

	version, err := DetermineSelectedPythonVersion(ReadState())

	if err != nil || version != mockState.GlobalVersion {
		t.Errorf("Expected version to be %s, got %s instead.", mockState.GlobalVersion, version)
	}
}

func TestDetermineSelectedPythonVersionDefaultsToSystem(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	version, err := DetermineSelectedPythonVersion(ReadState())

	if err != nil || version != "SYSTEM" {
		t.Errorf("Expected version to be 'SYSTEM', got %s instead.", version)
	}
}

func TestSearchForPythonVersionFileFindsFileInCwd(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	temporaryWd := t.TempDir()
	os.Chdir(temporaryWd)
	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)

	versionFound, found := SearchForPythonVersionFile()

	if versionFound != "1.2.3" || !found {
		t.Errorf("Expected \"1.2.3\", found %s", versionFound)
	}
}

func TestSearchForPythonVersionFileFindsFileInParents(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	temporaryWd := t.TempDir()

	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)
	os.Mkdir(path.Join(temporaryWd, "child"), 0750)
	os.Chdir(path.Join(temporaryWd, "child"))

	versionFound, found := SearchForPythonVersionFile()

	if versionFound != "1.2.3" || !found {
		t.Errorf("Expected \"1.2.3\", found %s", versionFound)
	}

}

func TestSearchForPythonVersionFileReturnsOnRootIfNoneFound(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	versionFound, found := SearchForPythonVersionFile()

	if versionFound != "" || found {
		t.Errorf("Did not expect any result, found %s.", versionFound)
	}

}
