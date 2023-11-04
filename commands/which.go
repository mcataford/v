package subcommands

import (
	"fmt"
	argparse "v/argparse"
	stateManager "v/state"
)

func Which(args []string, flags argparse.Flags, currentState stateManager.State) error {
	fmt.Println(currentState.GlobalVersion)
	return nil
}
