package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yv35.com/dotfiles-cli/internal/util/sh"
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

// GetRemoteURL returns the remote origin URL for the given repository
// Returns empty string if the repository has no remote or if there's an error
func GetRemoteURL(repoPath string) string {
	output, err := sh.RunShellOutput(fmt.Sprintf("git -C %s remote get-url origin", repoPath))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}
