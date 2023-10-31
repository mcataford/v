package subcommands

import (
	"fmt"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
)

func Which(args []string, flags argparse.Flags, currentState stateManager.State) error {
	fmt.Println(currentState.GlobalVersion)
	return nil
}
