package subcommands

import (
	"fmt"
	"os"
	"path"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
)

func UninstallPython(args []string, flags argparse.Flags, currentState stateManager.State) error {
	runtimePath := stateManager.GetPathFromStateDirectory(path.Join("runtimes", fmt.Sprintf("py-%s", args[1])))
	err := os.RemoveAll(runtimePath)
	return err
}
