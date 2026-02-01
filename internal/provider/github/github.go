package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/archive"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "github"
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no github packages defined
	if len(group.Github) == 0 {
		return nil
	}

	for _, githubSpec := range group.Github {
		if err := p.installGithubRelease(githubSpec, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) installGithubRelease(spec config.GithubSpec, onComplete types.OnTaskComplete) error {
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

	// Fetch release info from GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", spec.Repo)
	if spec.Version != "" && spec.Version != "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/v%s", spec.Repo, spec.Version)
	}

	resp, err := http.Get(url)
	if err != nil {
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to fetch release info for %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status %d", resp.StatusCode)
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to fetch release info for %s: %w", name, err)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to decode release info for %s: %w", name, err)
	}

	// Find matching asset
	re, err := regexp.Compile(spec.AssetPattern)
	if err != nil {
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("invalid asset pattern %s: %w", spec.AssetPattern, err)
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if re.MatchString(asset.Name) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		err := fmt.Errorf("no asset matching pattern %s found in release %s", spec.AssetPattern, release.TagName)
		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return err
	}

	// Download the asset
	tmpFile, err := archive.DownloadFile(downloadURL, name)
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
