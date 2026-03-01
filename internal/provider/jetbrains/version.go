package jetbrains

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// productInfo mirrors the relevant fields of product-info.json present in
// every modern JetBrains IDE installation.
type productInfo struct {
	Version     string `json:"version"`
	BuildNumber string `json:"buildNumber"`
}

// productInfoPath returns the path to product-info.json inside the install dir.
// On macOS the file lives inside the .app bundle at Contents/Resources/;
// on Linux it sits directly in the install dir root.
func productInfoPath(installDir string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(installDir, "Contents", "Resources", "product-info.json")
	}
	return filepath.Join(installDir, "product-info.json")
}

// getInstalledVersion reads the version from product-info.json.
// Returns an empty string (no error) when the file does not exist.
func getInstalledVersion(dir string) (string, error) {
	path := productInfoPath(dir)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var info productInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return "", err
	}

	return info.Version, nil
}
