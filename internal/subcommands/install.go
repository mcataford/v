package subcommands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	argparse "v/internal/argparse"
	stateManager "v/internal/state"
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

	fmt.Printf("🔨 Downloading source for Python %s\n", version)

	resp, err := client.Get(sourceUrl)

	if err != nil {
		return PackageMetadata{}, err
	}

	defer resp.Body.Close()

	io.Copy(file, resp.Body)

	defer file.Close()

	return PackageMetadata{ArchivePath: archivePath, Version: version}, nil
}

func buildFromSource(pkgMeta PackageMetadata, verbose bool) (PackageMetadata, error) {
	fmt.Printf("🔨 Unpacking source for %s\n", pkgMeta.ArchivePath)

	_, untarErr := RunCommand([]string{"tar", "zxvf", pkgMeta.ArchivePath}, stateManager.GetPathFromStateDirectory("cache"), !verbose)

	if untarErr != nil {
		return pkgMeta, untarErr
	}

	unzippedRoot := strings.TrimSuffix(pkgMeta.ArchivePath, path.Ext(pkgMeta.ArchivePath))

	fmt.Println("🔨 Building from source...")

	targetDirectory := stateManager.GetPathFromStateDirectory(path.Join("runtimes", fmt.Sprintf("py-%s", pkgMeta.Version)))

	_, configureErr := RunCommand([]string{"./configure", fmt.Sprintf("--prefix=%s", targetDirectory), "--enable-optimizations"}, unzippedRoot, !verbose)

	if configureErr != nil {
		return pkgMeta, configureErr
	}
	_, buildErr := RunCommand([]string{"make", "altinstall"}, unzippedRoot, !verbose)

	if buildErr != nil {
		return pkgMeta, buildErr
	}

	pkgMeta.InstallPath = targetDirectory

	fmt.Printf("Installed Python %s at %s\n", pkgMeta.Version, pkgMeta.InstallPath)
	fmt.Println("✅ All done!")

	return pkgMeta, nil
}
