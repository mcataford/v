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

func RunCommand(command []string, cwd string, quiet bool) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)

	cmd.Dir = cwd

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return "", nil
}
