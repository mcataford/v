package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Persistent state used by the CLI to track runtime information
// between calls.
type State struct {
	GlobalVersion string `json:"globalVersion"`
}

func GetStateDirectory() string {
	home, _ := os.UserHomeDir()
	userDefinedRoot, found := os.LookupEnv("V_ROOT")

	root := path.Join(home, ".v")

	if found {
		root = userDefinedRoot
	}

	return root
}

func GetPathFromStateDirectory(suffix string) string {
	return path.Join(GetStateDirectory(), suffix)
}

func ReadState() State {
	c, _ := ioutil.ReadFile(GetPathFromStateDirectory("state.json"))

	state := State{}

	json.Unmarshal(c, &state)

	return state
}

func WriteState(version string) {
	state := State{GlobalVersion: version}

	d, _ := json.Marshal(state)
	ioutil.WriteFile(GetPathFromStateDirectory("state.json"), d, 0750)
}

func GetAvailableVersions() []string {
	entries, _ := os.ReadDir(GetPathFromStateDirectory("runtimes"))

	versions := []string{}

	for _, d := range entries {
		versions = append(versions, strings.TrimPrefix(d.Name(), "py-"))
	}

	return versions
}
