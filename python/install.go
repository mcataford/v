package python

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
	exec "v/exec"
	logger "v/logger"
	state "v/state"
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

func (t VersionTag) MajorMinor() string {
	return t.Major + "." + t.Minor
}

// Installing new distribution happens in three stages:
// 1. Validating that the version number is of a valid format;
// 2. Downloading the source tarball;
// 3. Unzipping + building from source.
//
// The tarball is cached in the `cache` state directory and is reused
// if the same version is installed again later.
func InstallPythonDistribution(version string, noCache bool) error {
	if err := ValidateVersion(version); err != nil {
		return err
	}

	packageMetadata, dlerr := downloadSource(version, noCache)

	if dlerr != nil {
		return dlerr
	}

	if _, err := buildFromSource(packageMetadata); err != nil {
		return err
	}

	return nil
}

// Fetches the Python tarball for version <version> from python.org.
func downloadSource(version string, skipCache bool) (PackageMetadata, error) {
	archiveName := "Python-" + version + ".tgz"
	archivePath := state.GetStatePath("cache", archiveName)
	sourceUrl, _ := url.JoinPath(pythonReleasesBaseURL, version, archiveName)

	logger.InfoLogger.Println(logger.Bold("Downloading source for Python " + version))
	logger.InfoLogger.SetPrefix("  ")
	defer logger.InfoLogger.SetPrefix("")

	start := time.Now()

	if _, err := os.Stat(archivePath); errors.Is(err, os.ErrNotExist) || skipCache {
		logger.InfoLogger.Println("Fetching from " + sourceUrl)

		resp, err := http.Get(sourceUrl)

		if err != nil {
			return PackageMetadata{}, err
		}

		defer resp.Body.Close()
		file, _ := os.Create(archivePath)
		io.Copy(file, resp.Body)

		defer file.Close()
	} else {
		logger.InfoLogger.Println("Found in cache: " + archivePath)
	}

	logger.InfoLogger.Printf("✅ Done (%s)\n", time.Since(start))
	return PackageMetadata{ArchivePath: archivePath, Version: version}, nil
}

func buildFromSource(pkgMeta PackageMetadata) (PackageMetadata, error) {
	logger.InfoLogger.Println(logger.Bold("Building from source"))
	logger.InfoLogger.SetPrefix("  ")
	defer logger.InfoLogger.SetPrefix("")

	start := time.Now()

	logger.InfoLogger.Println("Unpacking source for " + pkgMeta.ArchivePath)

	if _, untarErr := exec.RunCommand([]string{"tar", "zxvf", pkgMeta.ArchivePath}, state.GetStatePath("cache")); untarErr != nil {
		return pkgMeta, untarErr
	}

	unzippedRoot := strings.TrimSuffix(pkgMeta.ArchivePath, path.Ext(pkgMeta.ArchivePath))

	logger.InfoLogger.Println("Configuring installer")

	if _, err := os.Stat(state.GetStatePath("runtimes", "python")); os.IsNotExist(err) {
		os.Mkdir(state.GetStatePath("runtimes", "python"), 0775)
	}

	targetDirectory := state.GetStatePath("runtimes", "python", pkgMeta.Version)

	if _, configureErr := exec.RunCommand([]string{"./configure", "--prefix=" + targetDirectory, "--enable-optimizations"}, unzippedRoot); configureErr != nil {
		return pkgMeta, configureErr
	}

	logger.InfoLogger.Println("Building")

	if _, buildErr := exec.RunCommand([]string{"make", "altinstall", "-j4"}, unzippedRoot); buildErr != nil {
		return pkgMeta, buildErr
	}

	if cleanupErr := os.RemoveAll(unzippedRoot); cleanupErr != nil {
		return pkgMeta, cleanupErr
	}

	pkgMeta.InstallPath = targetDirectory

	logger.InfoLogger.Printf("✅ Installed Python %s at %s (%s)\n", pkgMeta.Version, pkgMeta.InstallPath, time.Since(start))
	return pkgMeta, nil
}
