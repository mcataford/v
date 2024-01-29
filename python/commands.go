package python

import (
	"os"
	"slices"
	cli "v/cli"
	logger "v/logger"
	state "v/state"
)

func uninstallPython(args []string, flags cli.Flags, currentState state.State) error {
	runtimePath := state.GetStatePath("runtimes", "python", args[1])
	err := os.RemoveAll(runtimePath)
	return err
}

func installPython(args []string, flags cli.Flags, currentState state.State) error {
	version := args[1]

	return InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
}

func use(args []string, flags cli.Flags, currentState state.State) error {
	version := args[1]
	if err := ValidateVersion(version); err != nil {
		return err
	}

	availableVersions := state.GetAvailableVersions()

	found := false
	for _, v := range availableVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		logger.InfoLogger.Println("Version not installed. Installing it first.")
		InstallPythonDistribution(version, flags.NoCache, flags.Verbose)
	}

	state.WriteState(version)
	logger.InfoLogger.Printf("Now using Python %s\n", version)

	return nil
}

func listVersions(args []string, flags cli.Flags, currentState state.State) error {
	installedVersions, err := ListInstalledVersions()

	if err != nil {
		return err
	}

	if len(installedVersions) == 0 {
		logger.InfoLogger.Println("No versions installed!")
		return nil
	}

	for _, d := range installedVersions {
		logger.InfoLogger.Println(d)
	}

	return nil
}

// Which prints out the system path to the executable being used by `python`.
func which(args []string, flags cli.Flags, currentState state.State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	var printedPath string

	if selectedVersion.Source == "system" {
		_, sysPath := DetermineSystemPython()
		printedPath = sysPath + " (system)"
	} else if isInstalled {
		tag := VersionStringToStruct(selectedVersion.Version)
		printedPath = state.GetStatePath("runtimes", "python", selectedVersion.Version, "bin", "python"+tag.MajorMinor())
	} else {
		logger.InfoLogger.Printf("The desired version (%s) is not installed.\n", selectedVersion.Version)
		return nil
	}

	prefix := "Python path: "

	if flags.RawOutput {
		prefix = ""
	} else {
		printedPath = logger.Bold(printedPath)
	}

	logger.InfoLogger.Printf("%s%s\n", prefix, printedPath)
	return nil
}

// CurrentVersion (called via `v version`) outputs the currently selected version
// and what configures it. If the version is configured by a file, the file is returned
// under "source", if the system Python is used, "system" is returned as a source.
func currentVersion(args []string, flags cli.Flags, currentState state.State) error {
	selectedVersion, _ := DetermineSelectedPythonVersion(currentState)
	installedVersions, _ := ListInstalledVersions()
	isInstalled := slices.Contains(installedVersions, selectedVersion.Version)

	if !isInstalled {
		logger.InfoLogger.Println(logger.Bold(logger.Yellow("WARNING: This version is not installed.")))
	}

	if flags.RawOutput {
		logger.InfoLogger.Println(selectedVersion.Version)
		return nil
	}

	logger.InfoLogger.Printf("Python version: %s\nSource: %s\n", logger.Bold(selectedVersion.Version), logger.Bold(selectedVersion.Source))
	return nil
}
