package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartFmtGroup(label string) func(string) {
	fmt.Printf("\033[1m%s\033[0m\n", label)

	return func(message string) {
		fmt.Printf("    %s\n", message)
	}
}

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
