package state

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// If an override is provided via V_ROOT, it's used as state path.
func TestGetStatePathUsesEnvVRootOverride(t *testing.T) {
	MOCK_ROOT_PATH := "/overriden_path_to_state"
	os.Setenv("V_ROOT", MOCK_ROOT_PATH)
	defer os.Unsetenv("V_ROOT")
	if statePath := GetStatePath(); statePath != MOCK_ROOT_PATH {
		t.Errorf("Did not find expected state path %s, found %s instead.", MOCK_ROOT_PATH, statePath)
	}
}

// If no override is provided via V_ROOT, $HOME/.v is used.
func TestGetStatePathUsesHomeDefaultIfNoOverride(t *testing.T) {
	home, _ := os.UserHomeDir()

	expected := path.Join(home, ".v")
	if statePath := GetStatePath(); statePath != expected {
		t.Errorf("Did not find expected state path %s, found %s instead.", expected, statePath)
	}
}

// State is read from the state path and returned by ReadState.
func TestReadStateReadsStateAtPath(t *testing.T) {
	tempFolder, _ := os.MkdirTemp("", "v_test")
	defer os.RemoveAll(tempFolder)
	os.Setenv("V_ROOT", tempFolder)
	defer os.Unsetenv("V_ROOT")

	mockState := State{GlobalVersion: "1.0.0"}
	statePath := GetStatePath("state.json")
	stateData, _ := json.Marshal(mockState)
	ioutil.WriteFile(statePath, stateData, 0750)

	readState := ReadState()

	if readState != mockState {
		t.Errorf("Did not find expected state. %v != %v", mockState, readState)
	}
}

func TestWriteStateWritesAtPath(t *testing.T) {
	tempFolder, _ := os.MkdirTemp("", "v_test")
	defer os.RemoveAll(tempFolder)
	os.Setenv("V_ROOT", tempFolder)
	defer os.Unsetenv("V_ROOT")

	mockState := State{GlobalVersion: "2.0.0"}
	WriteState(mockState.GlobalVersion)

	statePath := GetStatePath("state.json")
	readState := State{}
	bytes, _ := ioutil.ReadFile(statePath)
	json.Unmarshal(bytes, &readState)

	if readState != mockState {
		t.Errorf("Did not find expected state. %v != %v", mockState, readState)
	}
}
