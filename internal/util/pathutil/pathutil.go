package pathutil

import (
	"os"
	"strings"
)

// MinimizePath replaces the user's home directory with ~ for shorter display
func MinimizePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	path = strings.ReplaceAll(path, homeDir, "~")
	return path
}

// MustHomeDir returns the user's home directory or panics if unavailable
func MustHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
