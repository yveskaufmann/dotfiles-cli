package fsutil

import (
	"fmt"
	"os"
)

// EnsureDirectory attempts to create a given directory at the specified directoryPath.
// In comparison to os.MkdirAll is it return application specific errors.
func EnsureDirectory(directoryPath string) error {
	info, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(directoryPath, 0755)
	}

	if err != nil {
		return fmt.Errorf("failed to access %s: %w ", directoryPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s exists but is not a directory", directoryPath)
	}

	return nil
}
