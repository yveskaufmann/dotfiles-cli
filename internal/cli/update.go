package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"yv35.com/dotfiles-cli/internal/util/archive"
)

const (
	updateRepo    = "yveskaufmann/dotfiles-cli"
	updateBinName = "dotfiles"
)

type selfUpdateRelease struct {
	TagName string `json:"tag_name"`
	Draft   bool   `json:"draft"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update dotfiles CLI to the latest release",
	Long:  `Check GitHub for the latest release and replace the current binary if a newer version is available.`,
	RunE:  runUpdate,
}

func init() {
	RegisterCommands(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking for updates...")

	release, err := fetchLatestRelease(updateRepo)
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}

	latestVersion := release.TagName
	currentVersion := versionInfo.version

	if !isNewerVersion(latestVersion, currentVersion) {
		fmt.Printf("✅ Already up to date (%s)\n", normalizeVersion(currentVersion))
		return nil
	}

	fmt.Printf("🆕 New version available: %s (current: %s)\n",
		normalizeVersion(latestVersion), normalizeVersion(currentVersion))

	assetURL, err := findUpdateAsset(release)
	if err != nil {
		return fmt.Errorf("no suitable asset found for %s/%s: %w", runtime.GOOS, runtime.GOARCH, err)
	}

	tmpFile, err := archive.DownloadFile(assetURL, updateBinName)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpFile)

	currentBin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine current binary path: %w", err)
	}
	targetDir := filepath.Dir(currentBin)

	if err := archive.InstallBinary(tmpFile, targetDir, updateBinName); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	fmt.Printf("🎉 Updated to %s\n", normalizeVersion(latestVersion))
	return nil
}

// fetchLatestRelease fetches the latest non-draft, non-prerelease release from GitHub.
func fetchLatestRelease(repo string) (*selfUpdateRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release selfUpdateRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// findUpdateAsset picks the asset matching the current OS and architecture.
// Asset naming convention: dotfiles-cli_{version}_{os}_{arch}.tar.gz  (linux/darwin)
//
//	dotfiles-cli_{version}_{os}_{arch}.zip   (windows)
func findUpdateAsset(release *selfUpdateRelease) (string, error) {
	version := strings.TrimPrefix(release.TagName, "v")
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	want := fmt.Sprintf("dotfiles-cli_%s_%s_%s.%s", version, runtime.GOOS, runtime.GOARCH, ext)

	for _, asset := range release.Assets {
		if asset.Name == want {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no asset named %q in release %s", want, release.TagName)
}

// isNewerVersion reports whether latest is strictly newer than current.
// A current version of "dev" is always treated as older so developer builds
// can exercise the update flow.
func isNewerVersion(latest, current string) bool {
	if current == "dev" {
		return true
	}
	lv := parseVersion(latest)
	cv := parseVersion(current)
	for i := 0; i < 3; i++ {
		if lv[i] > cv[i] {
			return true
		}
		if lv[i] < cv[i] {
			return false
		}
	}
	return false
}

// parseVersion parses a semver string (with optional "v" prefix) into [major, minor, patch].
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		// Strip any pre-release suffix (e.g. "1-beta" → 1)
		p, _, _ = strings.Cut(p, "-")
		n, _ := strconv.Atoi(p)
		result[i] = n
	}
	return result
}

// normalizeVersion ensures the version string has a "v" prefix for display.
func normalizeVersion(v string) string {
	if v == "dev" || strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}
