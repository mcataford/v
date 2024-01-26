package exec

import (
	"os"
	"os/exec"
	"strings"
)

// RunCommand is a thin wrapper around running command-line calls
// programmatically. It abstracts common configuration like routing
// output and handling the directory the calls are made from.
func RunCommand(command []string, cwd string, quiet bool) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)

	cmd.Dir = cwd

	var out strings.Builder
	var errOut strings.Builder

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = &out
		cmd.Stderr = &errOut
	}

	if err := cmd.Run(); err != nil {
		return errOut.String(), err
	}

	return out.String(), nil
}
