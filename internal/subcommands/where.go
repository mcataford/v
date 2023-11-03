package subcommands

import (
	"fmt"
	"strings"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
)

func versionStringToStruct(version string) VersionTag {
	splitVersion := strings.Split(version, ".")

	return VersionTag{Major: splitVersion[0], Minor: splitVersion[1], Patch: splitVersion[2]}
}

func Where(args []string, flags argparse.Flags, currentState stateManager.State) error {
	version := currentState.GlobalVersion
	tag := versionStringToStruct(version)
	withoutPatch := fmt.Sprintf("%s.%s", tag.Major, tag.Minor)
	fmt.Printf("%s/runtimes/py-%s/bin/python%s\n", stateManager.GetStateDirectory(), currentState.GlobalVersion, withoutPatch)
	return nil
}
