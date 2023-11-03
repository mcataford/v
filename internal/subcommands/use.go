package subcommands

import (
	"errors"
	"fmt"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
)

func Use(args []string, flags argparse.Flags, currentState stateManager.State) error {
	version := args[1]
	if err := validateVersion(version); err != nil {
		return err
	}

	availableVersions := stateManager.GetAvailableVersions()

	found := false
	for _, v := range availableVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		return errors.New("Version not installed.")
	}

	stateManager.WriteState(version)
	fmt.Printf("Now using Python %s\n", version)

	return nil
}
