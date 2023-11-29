package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func VersionStringToStruct(version string) VersionTag {
	splitVersion := strings.Split(version, ".")

	return VersionTag{Major: splitVersion[0], Minor: splitVersion[1], Patch: splitVersion[2]}
}

func ValidateVersion(version string) error {
	if splitVersion := strings.Split(version, "."); len(splitVersion) != 3 {
		return errors.New("Invalid version string. Expected format 'a.b.c'.")
	}

	return nil
}

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
