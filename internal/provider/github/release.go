package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"yv35.com/dotfiles-cli/internal/config"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// fetchReleaseInfo fetches GitHub release information from the API
func (p *Provider) fetchReleaseInfo(spec config.GithubSpec) (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", spec.Repo)
	if spec.Version != "" && spec.Version != "latest" {

		version := spec.Version
		if matched, _ := regexp.MatchString(`^\d+(\.\d+)*$`, spec.Version); matched {
			version = "v" + spec.Version
		}

		url = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", spec.Repo, version)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release info: status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// findAssetURL finds the download URL for a matching asset in the release
func (p *Provider) findAssetURL(release *githubRelease, assetPattern string) (string, error) {
	re, err := regexp.Compile(assetPattern)
	if err != nil {
		return "", fmt.Errorf("invalid asset pattern %s: %w", assetPattern, err)
	}

	for _, asset := range release.Assets {
		if re.MatchString(asset.Name) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no asset matching pattern %s found in release %s", assetPattern, release.TagName)
}
