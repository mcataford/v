package subcommands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	argparse "v/argparse"
	stateManager "v/state"
	util "v/util"
)

var pythonReleasesBaseURL = "https://www.python.org/ftp/python"

type PackageMetadata struct {
	ArchivePath string
	InstallPath string
	Version     string
}

type VersionTag struct {
	Major string
	Minor string
	Patch string
}

func InstallPython(args []string, flags argparse.Flags, currentState stateManager.State) error {
	verbose := flags.Verbose
	version := args[1]

	if err := validateVersion(version); err != nil {
		return err
	}

	packageMetadata, dlerr := downloadSource(version, "")

	if dlerr != nil {
		return dlerr
	}
	_, err := buildFromSource(packageMetadata, verbose)

	if err != nil {
		return err
	}

	return nil
}

// Fetches the Python tarball for version <version> from python.org
// and stores it at <destination>.
func downloadSource(version string, destination string) (PackageMetadata, error) {
	archiveName := fmt.Sprintf("Python-%s.tgz", version)
	archivePath := stateManager.GetPathFromStateDirectory(path.Join("cache", archiveName))
	sourceUrl := fmt.Sprintf("%s/%s/%s", pythonReleasesBaseURL, version, archiveName)
	file, _ := os.Create(archivePath)

	client := http.Client{}

	dlPrint := util.StartFmtGroup(fmt.Sprintf("Downloading source for Python %s", version))

	dlPrint(fmt.Sprintf("Fetching from %s", sourceUrl))
	start := time.Now()
	resp, err := client.Get(sourceUrl)

	if err != nil {
		return PackageMetadata{}, err
	}

	defer resp.Body.Close()

	io.Copy(file, resp.Body)

	defer file.Close()

	dlPrint(fmt.Sprintf("✅ Done (%s)", time.Since(start)))
	return PackageMetadata{ArchivePath: archivePath, Version: version}, nil
}

func buildFromSource(pkgMeta PackageMetadata, verbose bool) (PackageMetadata, error) {
	buildPrint := util.StartFmtGroup(fmt.Sprintf("Building from source"))
	start := time.Now()

	buildPrint(fmt.Sprintf("Unpacking source for %s", pkgMeta.ArchivePath))

	_, untarErr := RunCommand([]string{"tar", "zxvf", pkgMeta.ArchivePath}, stateManager.GetPathFromStateDirectory("cache"), !verbose)

	if untarErr != nil {
		return pkgMeta, untarErr
	}

	unzippedRoot := strings.TrimSuffix(pkgMeta.ArchivePath, path.Ext(pkgMeta.ArchivePath))

	buildPrint("Configuring installer")

	targetDirectory := stateManager.GetPathFromStateDirectory(path.Join("runtimes", fmt.Sprintf("py-%s", pkgMeta.Version)))

	_, configureErr := RunCommand([]string{"./configure", fmt.Sprintf("--prefix=%s", targetDirectory), "--enable-optimizations"}, unzippedRoot, !verbose)

	if configureErr != nil {
		return pkgMeta, configureErr
	}

	buildPrint("Building")
	_, buildErr := RunCommand([]string{"make", "altinstall"}, unzippedRoot, !verbose)

	if buildErr != nil {
		return pkgMeta, buildErr
	}

	if cleanupErr := os.RemoveAll(unzippedRoot); cleanupErr != nil {
		return pkgMeta, cleanupErr
	}

	pkgMeta.InstallPath = targetDirectory

	buildPrint(fmt.Sprintf("Installed Python %s at %s", pkgMeta.Version, pkgMeta.InstallPath))
	buildPrint(fmt.Sprintf("✅ Done (%s)", time.Since(start)))
	return pkgMeta, nil
}
