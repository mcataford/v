package subcommands

import (
	"fmt"
	"os"
	"strings"
	argparse "v/argparse"
	stateManager "v/state"
)

func ListVersions(args []string, flags argparse.Flags, currentState stateManager.State) error {
	runtimesDir := stateManager.GetPathFromStateDirectory("runtimes")
	entries, err := os.ReadDir(runtimesDir)

	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No versions installed!")
		return nil
	}

	for _, d := range entries {
		fmt.Println(strings.TrimPrefix(d.Name(), "py-"))
	}

	return nil
}
