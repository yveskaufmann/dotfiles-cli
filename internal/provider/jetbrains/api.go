package jetbrains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// IDECodes maps all known aliases (user-facing names and short codes) to the
// canonical JetBrains API product code.
var IDECodes = map[string]string{
	// Short codes (canonical)
	"IIU": "IIU",
	"IIC": "IIC",
	"PCP": "PCP",
	"PCC": "PCC",
	"PS":  "PS",
	"WS":  "WS",
	"CL":  "CL",
	"GO":  "GO",
	"RD":  "RD",
	"RM":  "RM",
	// Long-name aliases
	"idea-IU":                 "IIU",
	"idea":                    "IIC",
	"IntelliJ IDEA Ultimate":  "IIU",
	"IntelliJ IDEA Community": "IIC",
	"PyCharm Professional":    "PCP",
	"PyCharm Community":       "PCC",
	"PhpStorm":                "PS",
	"WebStorm":                "WS",
	"clion":                   "CL",
	"GoLand":                  "GO",
	"Rider":                   "RD",
	"RubyMine":                "RM",
}

// linuxInstallDirName maps an IDE product code to a directory name used under /opt on Linux.
var linuxInstallDirName = map[string]string{
	"IIU": "idea-IU",
	"IIC": "idea-IC",
	"PCP": "pycharm",
	"PCC": "pycharm-ce",
	"PS":  "phpstorm",
	"WS":  "webstorm",
	"CL":  "clion",
	"GO":  "goland",
	"RD":  "rider",
	"RM":  "rubymine",
}

// macAppName maps an IDE product code to its macOS .app bundle name.
var macAppName = map[string]string{
	"IIU": "IntelliJ IDEA.app",
	"IIC": "IntelliJ IDEA CE.app",
	"PCP": "PyCharm.app",
	"PCC": "PyCharm CE.app",
	"PS":  "PhpStorm.app",
	"WS":  "WebStorm.app",
	"CL":  "CLion.app",
	"GO":  "GoLand.app",
	"RD":  "Rider.app",
	"RM":  "RubyMine.app",
}

// resolveCode returns the canonical API code for an IDE identifier (code or
// alias). Returns an error when the identifier is unknown.
func resolveCode(ide string) (string, error) {
	if code, ok := IDECodes[ide]; ok {
		return code, nil
	}
	return "", fmt.Errorf("unknown JetBrains IDE identifier %q", ide)
}

// installDir returns the absolute installation path for an IDE code.
// On macOS it returns /Applications/<App>.app; on Linux /opt/<ide-dir>.
func installDir(code string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		app, ok := macAppName[code]
		if !ok {
			return "", fmt.Errorf("no macOS app name mapping for IDE code %q", code)
		}
		return "/Applications/" + app, nil
	default:
		dir, ok := linuxInstallDirName[code]
		if !ok {
			return "", fmt.Errorf("no install directory mapping for IDE code %q", code)
		}
		return "/opt/" + dir, nil
	}
}

// ReleaseInfo holds metadata for a single JetBrains release.
type ReleaseInfo struct {
	Version     string
	Build       string
	DownloadURL string
}

// jetbrainsReleasesURL is the JetBrains data service endpoint.
const jetbrainsReleasesURL = "https://data.services.jetbrains.com/products/releases"

// FetchRelease fetches the release info for the given IDE code from the
// JetBrains public API.  When version is empty or "latest", the latest stable
// release is returned; otherwise the exact version is matched.
func FetchRelease(ideCode, version string) (*ReleaseInfo, error) {
	latest := version == "" || strings.EqualFold(version, "latest")

	url := fmt.Sprintf(
		"%s?code=%s&type=release&_=%d",
		jetbrainsReleasesURL, ideCode, time.Now().UnixMilli(),
	)
	if latest {
		url += "&latest=true"
	}

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JetBrains releases for %s: %w", ideCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JetBrains API returned status %d for %s", resp.StatusCode, ideCode)
	}

	// The API returns: { "<CODE>": [ { "version": "...", "build": "...",
	// "downloads": { "linux": { "link": "...", ... }, ... } }, ... ] }
	var raw map[string][]struct {
		Version   string `json:"version"`
		Build     string `json:"build"`
		Downloads map[string]struct {
			Link         string `json:"link"`
			ChecksumLink string `json:"checksumLink"`
			Size         int64  `json:"size"`
		} `json:"downloads"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode JetBrains API response: %w", err)
	}

	releases, ok := raw[ideCode]
	if !ok || len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for IDE code %q", ideCode)
	}

	// When version is pinned, find the matching entry.
	var entry struct {
		Version   string
		Build     string
		Downloads map[string]struct {
			Link         string `json:"link"`
			ChecksumLink string `json:"checksumLink"`
			Size         int64  `json:"size"`
		}
	}
	found := false
	for _, r := range releases {
		if latest || r.Version == version {
			entry.Version = r.Version
			entry.Build = r.Build
			entry.Downloads = r.Downloads
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("version %q not found for IDE %s", version, ideCode)
	}

	// Pick the right download key for the current OS and prefer .tar.gz.
	downloadURL, err := pickDownloadURL(entry.Downloads)
	if err != nil {
		return nil, fmt.Errorf("IDE %s %s: %w", ideCode, entry.Version, err)
	}

	return &ReleaseInfo{
		Version:     entry.Version,
		Build:       entry.Build,
		DownloadURL: downloadURL,
	}, nil
}

// pickDownloadURL selects the most appropriate download URL for the current OS/arch.
// On macOS it returns a .dmg; on Linux it returns a .tar.gz.
func pickDownloadURL(downloads map[string]struct {
	Link         string `json:"link"`
	ChecksumLink string `json:"checksumLink"`
	Size         int64  `json:"size"`
}) (string, error) {
	type candidate struct {
		key    string
		suffix string
	}
	var candidates []candidate

	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			candidates = []candidate{{"macM1", ".dmg"}, {"mac", ".dmg"}}
		} else {
			candidates = []candidate{{"mac", ".dmg"}, {"macM1", ".dmg"}}
		}
	case "linux":
		if runtime.GOARCH == "arm64" {
			candidates = []candidate{{"linuxARM64", ".tar.gz"}, {"linux", ".tar.gz"}}
		} else {
			candidates = []candidate{{"linux", ".tar.gz"}, {"linuxARM64", ".tar.gz"}}
		}
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	for _, c := range candidates {
		d, ok := downloads[c.key]
		if !ok || d.Link == "" {
			continue
		}
		if !strings.HasSuffix(d.Link, c.suffix) {
			continue
		}
		return d.Link, nil
	}

	available := make([]string, 0, len(downloads))
	for k, d := range downloads {
		available = append(available, fmt.Sprintf("%s=%s", k, d.Link))
	}
	return "", fmt.Errorf(
		"no suitable download found for OS %s/%s (available: %s)",
		runtime.GOOS, runtime.GOARCH, strings.Join(available, ", "),
	)
}
