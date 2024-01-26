package cli

import (
	"testing"
	state "v/state"
)

func TestNamespaceAddCommand(t *testing.T) {
	namespace := Namespace{}

	canary := 0

	handler := func(a []string, b Flags, c state.State) error {
		canary = 1
		return nil
	}

	namespace.AddCommand("test", handler, "", "")

	if len(namespace.Commands) != 1 {
		t.Errorf("Expected one command, found %d", len(namespace.Commands))
	}

	if e := namespace.Commands["test"].Handler([]string{}, Flags{}, state.State{}); e != nil {
		t.Errorf("Unexpected error when running handler: %s", e)
	}

	if canary != 1 {
		t.Errorf("Expected canary value to have been modified.")
	}
}
