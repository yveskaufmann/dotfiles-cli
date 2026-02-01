package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"yv35.com/dotfiles/internal/tool"
	"yv35.com/dotfiles/internal/util/sh"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// InstallGithubRelease installs a binary from a GitHub release.
func InstallGithubRelease(name, repo, version, pattern, binaryPath, installPath string, binaries []string) error {
	// Check if already installed
	if len(binaries) > 0 {
		allInstalled := true
		for _, b := range binaries {
			if err := sh.RunShell("type " + b + " > /dev/null 2>&1"); err != nil {
				allInstalled = false
				break
			}
		}
		if allInstalled {
			fmt.Printf("✅ Binaries %v are already installed\n", binaries)
			return nil
		}
	} else if err := sh.RunShell("type " + name + " > /dev/null 2>&1"); err == nil {
		fmt.Printf("✅ %s is already installed (type check passed)\n", name)
		return nil
	}

	fmt.Printf("⚙️  Installing %s from GitHub: %s\n", name, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	if version != "" && version != "latest" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/v%s", repo, version)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch release info: status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	var downloadURL string
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid asset pattern %s: %w", pattern, err)
	}

	for _, asset := range release.Assets {
		if re.MatchString(asset.Name) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no asset matching pattern %s found in release %s", pattern, release.TagName)
	}

	tmpFile := filepath.Join(os.TempDir(), filepath.Base(downloadURL))
	defer os.Remove(tmpFile)

	fmt.Printf("⬇️  Downloading %s...\n", downloadURL)
	if err := sh.RunShell(fmt.Sprintf("curl -L -sS -o %s %s", tmpFile, downloadURL)); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	if len(binaries) > 0 {
		return tool.InstallMultipleFromArchive(binaries, tmpFile, installPath)
	}

	return tool.InstallFromArchiveOrBinary(name, tmpFile, binaryPath, installPath)
}
