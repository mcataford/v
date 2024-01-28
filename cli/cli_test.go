package cli

import (
	"slices"
	"testing"
)

func TestAddNamespace(t *testing.T) {
	cli := CLI{}

	namespace := Namespace{Label: "test"}

	cli.AddNamespace(namespace)

	if len(cli.Namespaces) != 1 {
		t.Errorf("Expected one namespace added.")
	}

	if cli.Namespaces["test"].Label != namespace.Label {
		t.Errorf("Unexpected label value: %s", cli.Namespaces["test"].Label)
	}
}

func TestListNamespacesReturnsAlphaOrderedLabels(t *testing.T) {
	cli := CLI{}

	n1 := Namespace{Label: "a"}
	n2 := Namespace{Label: "b"}
	cli.AddNamespace(n2).AddNamespace(n1)

	labels := cli.ListNamespaces()

	if !slices.Equal(labels, []string{"a", "b"}) {
		t.Errorf("Expected labels to be alpha ordered, got %v", labels)
	}

}
