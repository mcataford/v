package testutils

import (
	"os"
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
