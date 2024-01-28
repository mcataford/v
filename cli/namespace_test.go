package cli

import (
	"slices"
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

func TestNamespaceListCommandsReturnsAlphaSortedLabels(t *testing.T) {
	namespace := Namespace{}

	handler := func(a []string, b Flags, c state.State) error {
		return nil
	}

	// Inserted in non-alpha order.
	namespace.AddCommand("b", handler, "", "")
	namespace.AddCommand("a", handler, "", "")

	labels := namespace.ListCommands()

	if !slices.Equal(labels, []string{"a", "b"}) {
		t.Errorf("Expected labels to be alpha-ordered. Got %v", labels)
	}
}
