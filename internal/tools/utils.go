package tools

import (
	"os"
	"os/exec"
	"path/filepath"
)

func IsBinaryOnPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunScript(path string) error {
	// If path is relative, try to resolve it relative to $DOTFILES
	if !filepath.IsAbs(path) {
		dotfiles := os.Getenv("DOTFILES")
		if dotfiles != "" {
			fullPath := filepath.Join(dotfiles, path)
			if _, err := os.Stat(fullPath); err == nil {
				path = fullPath
			}
		}
	}
	return RunShell(path)
}
