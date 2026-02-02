package binary

import (
	"fmt"
	"os"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/util/archive"
	"yv35.com/dotfiles/internal/util/stringutils"
)

// installBinary downloads and installs a binary from a URL
func (p *Provider) installBinary(spec config.BinarySpec) error {
	name := spec.Name

	// Replace version placeholder
	version := spec.Version
	if version == "" {
		version = "latest"
	}

	finalURL := stringutils.ResolvePlaceholdersWithVars(spec.URL, map[string]string{
		"version": version,
	})

	// Download the file
	tmpFile, err := archive.DownloadFile(finalURL, name)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", name, err)
	}
	defer os.Remove(tmpFile)

	// Install based on whether it's multiple binaries or single
	if len(spec.Binaries) > 0 {
		return p.installMultipleBinaries(spec.Binaries, tmpFile, spec.InstallPath)
	}
	return p.installSingleBinary(name, tmpFile, spec.BinaryPath, spec.InstallPath)
}

// installMultipleBinaries extracts and installs multiple binaries from an archive
func (p *Provider) installMultipleBinaries(binaries []string, srcFile, installPath string) error {
	// Extract archive
	extractDir := "/tmp/multi_extract_" + fmt.Sprintf("%d", os.Getpid())
	os.RemoveAll(extractDir)
	defer os.RemoveAll(extractDir)

	if err := archive.ExtractArchive(srcFile, extractDir); err != nil {
		return err
	}

	// Find and install each binary
	for _, name := range binaries {
		srcBinary, err := archive.FindBinary(extractDir, name)
		if err != nil {
			return fmt.Errorf("binary %s not found in archive: %w", name, err)
		}

		if err := archive.InstallBinary(srcBinary, installPath, name); err != nil {
			return err
		}
	}

	return nil
}

// installSingleBinary installs a single binary, handling both archives and direct binaries
func (p *Provider) installSingleBinary(name, srcFile, binaryPath, installPath string) error {
	// Check if it's an archive
	if !archive.IsArchive(srcFile) {
		// Direct binary installation
		return archive.InstallBinary(srcFile, installPath, name)
	}

	// Extract archive
	extractDir := "/tmp/" + name + "_extract_" + fmt.Sprintf("%d", os.Getpid())
	os.RemoveAll(extractDir)
	defer os.RemoveAll(extractDir)

	if err := archive.ExtractArchive(srcFile, extractDir); err != nil {
		return err
	}

	// Find the binary
	srcBinary := ""
	if binaryPath != "" {
		srcBinary = extractDir + "/" + binaryPath
	} else {
		var err error
		srcBinary, err = archive.FindBinary(extractDir, name)
		if err != nil {
			return err
		}
	}

	return archive.InstallBinary(srcBinary, installPath, name)
}
