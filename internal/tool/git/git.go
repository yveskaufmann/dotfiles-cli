package git

import (
	"fmt"

	"yv35.com/dotfiles/internal/tool"
	"yv35.com/dotfiles/internal/util/sh"
)

// Ensure will attempt to ensure that git is installed on the system.
// If git is not installed, it will try to install it using the system's package manager,
// when it fails it will return an error.
func Ensure() error {

	if sh.IsBinaryOnPath("git") {
		fmt.Printf("✅ git is already installed\n")
		return nil
	}

	fmt.Printf("⚙️  git is not installed, attempting to install...\n")

	if err := tool.InstallAptPackages([]string{"git"}); err != nil {
		return fmt.Errorf("failed to install git: %w", err)
	}

	fmt.Printf("✅ git has been installed successfully\n")

	return nil
}
