package main

import (
	"bytes"
	"os"
	"testing"
)

func TestListVersionOutputsNoticeIfNoVersionsInstalled(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	os.Mkdir(GetStatePath("runtimes"), 0750)
	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	ListVersions([]string{}, Flags{}, State{})

	captured := out.String()
	if captured != "No versions installed!\n" {
		t.Errorf("Unexpected message: %s", captured)
	}
}

func TestListVersionOutputsVersionsInstalled(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	os.MkdirAll(GetStatePath("runtimes", "py-1.2.3"), 0750)
	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	ListVersions([]string{}, Flags{}, State{})

	captured := out.String()
	if captured != "1.2.3\n" {
		t.Errorf("Unexpected message: %s", captured)
	}
}

func TestListVersionReturnsErrorOnFailure(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	err := ListVersions([]string{}, Flags{}, State{})

	captured := out.String()
	if captured != "" {
		t.Errorf("Captured unexpected message: %s", captured)
	}

	if err == nil {
		t.Errorf("Expected error returned, did not get one.")
	}
}
