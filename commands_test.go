package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	state "v/state"
	testutils "v/testutils"
)

func TestWriteShim(t *testing.T) {
	defer testutils.SetupAndCleanupEnvironment(t)()

	os.Mkdir(state.GetStatePath("shims"), 0775)
	testShimPath := state.GetStatePath("shims", "testshim")
	e := writeShim(testShimPath)

	shimContent, _ := ioutil.ReadFile(testShimPath)

	if e != nil {
		t.Errorf("Errored while writing shim")
	}

	if !strings.Contains(string(shimContent), "$(v python which --raw) $@") {
		t.Errorf("Expected shim to contain pass-through via 'which', got %s", shimContent)
	}

}
