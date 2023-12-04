package main

import (
	"bytes"
	"os"
	"strings"
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

func TestListVersionOutputsVersionSelectedAndWarnsNotInstalled(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	Which([]string{}, Flags{}, State{GlobalVersion: "1.2.3"})

	captured := out.String()
	if captured != "The desired version (1.2.3) is not installed.\n" {
		t.Errorf("Unexpected message: %s", captured)
	}
}

func TestWhichOutputsVersionSelectedIfInstalled(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	os.MkdirAll(GetStatePath("runtimes", "py-1.2.3"), 0750)
	Which([]string{}, Flags{}, State{GlobalVersion: "1.2.3"})

	captured := strings.TrimSpace(out.String())
	expected := GetStatePath("runtimes", "py-1.2.3", "bin", "python1.2")
	if !strings.Contains(captured, expected) {
		t.Errorf("Unexpected message: %s, not %s", captured, expected)
	}
}

func TestWhichOutputsSystemVersionIfNoneSelected(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	Which([]string{}, Flags{RawOutput: true}, State{})

	captured := strings.TrimSpace(out.String())

	if captured != "/bin/python (system)" {
		t.Errorf("%s != %s", captured, "/bin/python (system)")
	}
}

func TestWhichOutputsVersionWithoutPrefixesIfRawOutput(t *testing.T) {
	defer SetupAndCleanupEnvironment(t)()

	var out bytes.Buffer

	InfoLogger.SetOutput(&out)
	defer InfoLogger.SetOutput(os.Stdout)

	os.MkdirAll(GetStatePath("runtimes", "py-1.2.3"), 0750)
	Which([]string{}, Flags{RawOutput: true}, State{GlobalVersion: "1.2.3"})

	captured := strings.TrimSpace(out.String())
	expected := GetStatePath("runtimes", "py-1.2.3", "bin", "python1.2")
	if captured != expected {
		t.Errorf("Unexpected message: %s, not %s", captured, expected)
	}
}
