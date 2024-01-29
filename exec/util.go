package exec

import (
	"io"
	"os/exec"
	"strings"
	logger "v/logger"
)

// RunCommand is a thin wrapper around running command-line calls
// programmatically. It abstracts common configuration like routing
// output and handling the directory the calls are made from.
func RunCommand(command []string, cwd string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)

	cmd.Dir = cwd

	var out strings.Builder
	var errOut strings.Builder

	stdOutMultiWriter := io.MultiWriter(&out, logger.DebugLogger.Writer())
	stdErrMultiWriter := io.MultiWriter(&errOut, logger.DebugLogger.Writer())

	cmd.Stdout = stdOutMultiWriter
	cmd.Stderr = stdErrMultiWriter

	if err := cmd.Run(); err != nil {
		return errOut.String(), err
	}

	return out.String(), nil
}
