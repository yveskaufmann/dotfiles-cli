package git

import (
	"os"
	"path/filepath"

	"yv35.com/dotfiles/internal/util/sh"
)

// Clone clones a git repository from the given URL to the destination path
// Let git handle URL validation - any errors will be returned from git directly
func Clone(url, destination string) error {
	return sh.Run("git", "clone", url, destination)
}

// Pull pulls the latest changes from the remote repository
// Uses -C flag to run git command in the specified repository path
func Pull(repoPath string) error {
	return sh.Run("git", "-C", repoPath, "pull")
}

// IsRepository checks if the given path is a git repository
// by checking for the existence of a .git directory
func IsRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
