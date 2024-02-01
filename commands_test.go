package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	cli "v/cli"
	logger "v/logger"
	python "v/python"
	state "v/state"
	testutils "v/testutils"
)

func TestWriteShim(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	os.Mkdir(state.GetStatePath("shims"), 0775)
	testShimPath := state.GetStatePath("shims", "testshim")
	e := writeShim(testShimPath, "testcommand")

	shimContent, _ := ioutil.ReadFile(testShimPath)

	if e != nil {
		t.Errorf("Errored while writing shim")
	}

	if !strings.Contains(string(shimContent), "testcommand") {
		t.Errorf("Expected shim to contain pass-through via 'which', got %s", shimContent)
	}

}

func TestInitializeCreatesStateDirectories(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	err := Initialize([]string{}, cli.Flags{}, state.State{})

	if err != nil {
		t.Errorf("Unexpected error initializing")
	}

	if _, err = os.Stat(state.GetStatePath()); os.IsNotExist(err) {
		t.Errorf("Root state directory not found")
	}

	if _, err = os.Stat(state.GetStatePath("shims")); os.IsNotExist(err) {
		t.Errorf("Shims directory not found")
	}

	if _, err = os.Stat(state.GetStatePath("cache")); os.IsNotExist(err) {
		t.Errorf("Cache directory not found")
	}

	if _, err = os.Stat(state.GetStatePath("runtimes")); os.IsNotExist(err) {
		t.Errorf("Runtimes directory not found")
	}
}

func TestInitializeCreatesAllPythonShims(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	err := Initialize([]string{}, cli.Flags{}, state.State{})

	if err != nil {
		t.Errorf("Unexpected error initializing")
	}

	expectedShims := python.Shims

	for shimLabel, shimCall := range expectedShims {
		shimContent, err := os.ReadFile(state.GetStatePath("shims", shimLabel))

		if os.IsNotExist(err) {
			t.Errorf("%s shim not created", shimLabel)
		}

		if !strings.Contains(string(shimContent), shimCall) {
			t.Errorf("%s shim does not contain expected call (%s not in %s)", shimLabel, shimCall, shimContent)
		}
	}
}

func TestInitializeWithAddPathPrintsExportPATH(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	var buf bytes.Buffer

	logger.InfoLogger.SetOutput(&buf)

	defer func() {
		logger.InfoLogger.SetOutput(os.Stdout)
	}()

	err := Initialize([]string{}, cli.Flags{AddPath: true}, state.State{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "export PATH=" + state.GetStatePath("shims") + ":$PATH\n"
	if buf.String() != expected {
		t.Errorf("Expected PATH export, got %s", buf.String())
	}
}
