package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
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

func InstallPythonDistribution(version string, noCache bool, verbose bool) error {
	if err := ValidateVersion(version); err != nil {
		return err
	}

	packageMetadata, dlerr := downloadSource(version, noCache)

	if dlerr != nil {
		return dlerr
	}
	_, err := buildFromSource(packageMetadata, verbose)

	if err != nil {
		return err
	}

	return nil
}

// Fetches the Python tarball for version <version> from python.org.
func downloadSource(version string, skipCache bool) (PackageMetadata, error) {
	archiveName := "Python-" + version + ".tgz"
	archivePath := state.GetStatePath("cache", archiveName)
	sourceUrl, _ := url.JoinPath(pythonReleasesBaseURL, version, archiveName)

	client := http.Client{}

	logger.InfoLogger.Println(Bold("Downloading source for Python " + version))
	logger.InfoLogger.SetPrefix("  ")
	defer logger.InfoLogger.SetPrefix("")

	start := time.Now()
	_, err := os.Stat(archivePath)

	if errors.Is(err, os.ErrNotExist) || skipCache {
		logger.InfoLogger.Println("Fetching from " + sourceUrl)

		resp, err := client.Get(sourceUrl)

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

func buildFromSource(pkgMeta PackageMetadata, verbose bool) (PackageMetadata, error) {
	logger.InfoLogger.Println(Bold("Building from source"))
	logger.InfoLogger.SetPrefix("  ")
	defer logger.InfoLogger.SetPrefix("")

	start := time.Now()

	logger.InfoLogger.Println("Unpacking source for " + pkgMeta.ArchivePath)

	_, untarErr := RunCommand([]string{"tar", "zxvf", pkgMeta.ArchivePath}, state.GetStatePath("cache"), !verbose)

	if untarErr != nil {
		return pkgMeta, untarErr
	}

	unzippedRoot := strings.TrimSuffix(pkgMeta.ArchivePath, path.Ext(pkgMeta.ArchivePath))

	logger.InfoLogger.Println("Configuring installer")

	targetDirectory := state.GetStatePath("runtimes", "py-"+pkgMeta.Version)

	_, configureErr := RunCommand([]string{"./configure", "--prefix=" + targetDirectory, "--enable-optimizations"}, unzippedRoot, !verbose)

	if configureErr != nil {
		return pkgMeta, configureErr
	}

	logger.InfoLogger.Println("Building")
	_, buildErr := RunCommand([]string{"make", "altinstall", "-j4"}, unzippedRoot, !verbose)

	if buildErr != nil {
		return pkgMeta, buildErr
	}

	if cleanupErr := os.RemoveAll(unzippedRoot); cleanupErr != nil {
		return pkgMeta, cleanupErr
	}

	pkgMeta.InstallPath = targetDirectory

	logger.InfoLogger.Printf("Installed Python %s at %s\n", pkgMeta.Version, pkgMeta.InstallPath)
	logger.InfoLogger.Printf("✅ Done (%s)\n", time.Since(start))
	return pkgMeta, nil
}
