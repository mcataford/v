package python

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"slices"
	"testing"
	state "v/state"
	testutils "v/testutils"
)

func TestDetermineSystemPythonGetsUnshimmedPythonRuntime(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	ioutil.WriteFile(state.GetStatePath("shims", "python"), []byte("#!/bin/bash\necho \"Python 4.5.6\""), 0777)
	mockSystemPythonPath := t.TempDir()
	mockSystemPythonExecPath := path.Join(mockSystemPythonPath, "python")

	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:/usr/bin", state.GetStatePath("shims")))
	defer os.Setenv("PATH", oldPath)
	sysVersion, sysPath := DetermineSystemPython()

	if sysVersion == "4.5.6" {
		t.Errorf("Expected system Python to not match the shim, found %s instead.", sysVersion)
	}

	if sysPath != "/bin/python" {
		t.Errorf("Expected system Python path to be %s, found %s instead.", mockSystemPythonExecPath, sysPath)
	}
}

func TestDetermineSelectedPythonVersionUsesPythonVersionFileIfFound(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	// Writing a mock user-defined state.
	mockState := state.State{GlobalVersion: "1.0.0"}
	statePath := state.GetStatePath("state.json")
	stateData, _ := json.Marshal(mockState)
	ioutil.WriteFile(statePath, stateData, 0750)

	temporaryWd := t.TempDir()
	os.Chdir(temporaryWd)
	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)

	version, err := DetermineSelectedPythonVersion(state.ReadState())

	if err != nil || version.Version != "1.2.3" {
		t.Errorf("Expected version to be %s, got %s instead.", "1.2.3", version.Version)
	}
}

func TestDetermineSelectedPythonVersionGetsUserDefinedVersion(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	// Writing a mock user-defined state.
	mockState := state.State{GlobalVersion: "1.0.0"}
	statePath := state.GetStatePath("state.json")
	stateData, _ := json.Marshal(mockState)
	ioutil.WriteFile(statePath, stateData, 0750)

	version, err := DetermineSelectedPythonVersion(state.ReadState())

	if err != nil || version.Version != mockState.GlobalVersion {
		t.Errorf("Expected version to be %s, got %s instead.", mockState.GlobalVersion, version)
	}
}

func TestDetermineSelectedPythonVersionDefaultsToSystem(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	version, err := DetermineSelectedPythonVersion(state.ReadState())

	if err != nil || version.Source != "system" {
		t.Errorf("Expected version to be 'SYSTEM', got %s instead.", version)
	}
}

func TestSearchForPythonVersionFileFindsFileInCwd(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	temporaryWd := t.TempDir()
	os.Chdir(temporaryWd)
	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)

	versionFound, found := SearchForPythonVersionFile()

	if versionFound.Version != "1.2.3" || !found {
		t.Errorf("Expected \"1.2.3\", found %s", versionFound)
	}
}

func TestSearchForPythonVersionFileFindsFileInParents(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	temporaryWd := t.TempDir()

	ioutil.WriteFile(path.Join(temporaryWd, ".python-version"), []byte("1.2.3"), 0750)
	os.Mkdir(path.Join(temporaryWd, "child"), 0750)
	os.Chdir(path.Join(temporaryWd, "child"))

	versionFound, found := SearchForPythonVersionFile()

	if versionFound.Version != "1.2.3" || !found {
		t.Errorf("Expected \"1.2.3\", found %s", versionFound)
	}

}

func TestSearchForPythonVersionFileReturnsOnRootIfNoneFound(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	versionFound, found := SearchForPythonVersionFile()

	if versionFound.Version != "" || found {
		t.Errorf("Did not expect any result, found %s.", versionFound.Version)
	}
}

func TestListInstalledVersion(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	versions := []string{"1.2.3", "4.5.6", "7.8.9"}

	os.Mkdir(state.GetStatePath("runtimes"), 0750)
	for _, version := range versions {
		os.MkdirAll(state.GetStatePath("runtimes", "python", version), 0750)
	}

	installedVersions, _ := ListInstalledVersions()

	if !slices.Equal(installedVersions, versions) {
		t.Errorf("Expected %s, got %s.", versions, installedVersions)
	}
}

func TestListInstalledVersionNoVersionsInstalled(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	os.MkdirAll(state.GetStatePath("runtimes", "python"), 0750)

	installedVersions, _ := ListInstalledVersions()

	if len(installedVersions) != 0 {
		t.Errorf("Expected 0 elements, got %d (%s).", len(installedVersions), installedVersions)
	}
}

func TestListInstalledVersionNoRuntimesDir(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	installedVersions, err := ListInstalledVersions()

	if len(installedVersions) != 0 {
		t.Errorf("Expected 0 elements, got %d (%s).", len(installedVersions), installedVersions)
	}

	if err == nil {
		t.Errorf("Expected error to be returned, got nil.")
	}
}
