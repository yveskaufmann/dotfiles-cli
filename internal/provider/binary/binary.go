package binary

import (
	"fmt"
	"os"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/archive"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "binary"
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no binary packages defined
	if len(group.Binary) == 0 {
		return nil
	}

	for _, binarySpec := range group.Binary {
		if err := p.installBinary(binarySpec, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) installBinary(spec config.BinarySpec, onComplete types.OnTaskComplete) error {
	name := spec.Name

	// Check if already installed
	if len(spec.Binaries) > 0 {
		allInstalled := true
		for _, b := range spec.Binaries {
			if err := sh.RunShell("type " + b + " > /dev/null 2>&1"); err != nil {
				allInstalled = false
				break
			}
		}
		if allInstalled {
			onComplete(types.TaskResult{
				Name:   name,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("all binaries already installed"),
			})
			return nil
		}
	} else if err := sh.RunShell("type " + name + " > /dev/null 2>&1"); err == nil {
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusSkipped,
			Error:  fmt.Errorf("already installed"),
		})
		return nil
	}

	// Replace version placeholder
	version := spec.Version
	if version == "" {
		version = "latest"
	}
	finalURL := strings.ReplaceAll(spec.URL, ":version", version)

	// Download the file
	tmpFile, err := archive.DownloadFile(finalURL, name)
	if err != nil {
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to download %s: %w", name, err)
	}
	defer os.Remove(tmpFile)

	// Install based on whether it's multiple binaries or single
	if len(spec.Binaries) > 0 {
		if err := p.installMultipleBinaries(spec.Binaries, tmpFile, spec.InstallPath, onComplete); err != nil {
			onComplete(types.TaskResult{
				Name:   name,
				Status: types.StatusFailed,
				Error:  err,
			})
			return err
		}
	} else {
		if err := p.installSingleBinary(name, tmpFile, spec.BinaryPath, spec.InstallPath, onComplete); err != nil {
			onComplete(types.TaskResult{
				Name:   name,
				Status: types.StatusFailed,
				Error:  err,
			})
			return err
		}
	}

	onComplete(types.TaskResult{
		Name:   name,
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return nil
}

func (p *Provider) installMultipleBinaries(binaries []string, srcFile, installPath string, onComplete types.OnTaskComplete) error {
	// Extract archive
	extractDir := "/tmp/multi_extract"
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

func (p *Provider) installSingleBinary(name, srcFile, binaryPath, installPath string, onComplete types.OnTaskComplete) error {
	// Check if it's an archive
	if !archive.IsArchive(srcFile) {
		// Direct binary installation
		return archive.InstallBinary(srcFile, installPath, name)
	}

	// Extract archive
	extractDir := "/tmp/" + name + "_extract"
	os.RemoveAll(extractDir)
	defer os.RemoveAll(extractDir)

	if err := archive.ExtractArchive(srcFile, extractDir); err != nil {
		return err
	}

	// Find binary
	var srcBinary string
	if binaryPath != "" {
		srcBinary = extractDir + "/" + binaryPath
	} else {
		var err error
		srcBinary, err = archive.FindBinary(extractDir, name)
		if err != nil {
			return fmt.Errorf("binary %s not found in archive: %w", name, err)
		}
	}

	// Install binary
	return archive.InstallBinary(srcBinary, installPath, name)
}
