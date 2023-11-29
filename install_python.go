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
	archivePath := GetStatePath("cache", archiveName)
	sourceUrl, _ := url.JoinPath(pythonReleasesBaseURL, version, archiveName)

	client := http.Client{}

	InfoLogger.Println(Bold("Downloading source for Python " + version))
	InfoLogger.SetPrefix("  ")
	defer InfoLogger.SetPrefix("")

	start := time.Now()
	_, err := os.Stat(archivePath)

	if errors.Is(err, os.ErrNotExist) || skipCache {
		InfoLogger.Println("Fetching from " + sourceUrl)

		resp, err := client.Get(sourceUrl)

		if err != nil {
			return PackageMetadata{}, err
		}

		defer resp.Body.Close()
		file, _ := os.Create(archivePath)
		io.Copy(file, resp.Body)

		defer file.Close()
	} else {
		InfoLogger.Println("Found in cache: " + archivePath)
	}

	InfoLogger.Printf("✅ Done (%s)\n", time.Since(start))
	return PackageMetadata{ArchivePath: archivePath, Version: version}, nil
}

func buildFromSource(pkgMeta PackageMetadata, verbose bool) (PackageMetadata, error) {
	InfoLogger.Println(Bold("Building from source"))
	InfoLogger.SetPrefix("  ")
	defer InfoLogger.SetPrefix("")

	start := time.Now()

	InfoLogger.Println("Unpacking source for " + pkgMeta.ArchivePath)

	_, untarErr := RunCommand([]string{"tar", "zxvf", pkgMeta.ArchivePath}, GetStatePath("cache"), !verbose)

	if untarErr != nil {
		return pkgMeta, untarErr
	}

	unzippedRoot := strings.TrimSuffix(pkgMeta.ArchivePath, path.Ext(pkgMeta.ArchivePath))

	InfoLogger.Println("Configuring installer")

	targetDirectory := GetStatePath("runtimes", "py-"+pkgMeta.Version)

	_, configureErr := RunCommand([]string{"./configure", "--prefix=" + targetDirectory, "--enable-optimizations"}, unzippedRoot, !verbose)

	if configureErr != nil {
		return pkgMeta, configureErr
	}

	InfoLogger.Println("Building")
	_, buildErr := RunCommand([]string{"make", "altinstall", "-j4"}, unzippedRoot, !verbose)

	if buildErr != nil {
		return pkgMeta, buildErr
	}

	if cleanupErr := os.RemoveAll(unzippedRoot); cleanupErr != nil {
		return pkgMeta, cleanupErr
	}

	pkgMeta.InstallPath = targetDirectory

	InfoLogger.Printf("Installed Python %s at %s\n", pkgMeta.Version, pkgMeta.InstallPath)
	InfoLogger.Printf("✅ Done (%s)\n", time.Since(start))
	return pkgMeta, nil
}
