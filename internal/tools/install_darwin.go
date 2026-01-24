//go:build darwin

package tools

import "fmt"

func Install(packageName string) error {
	if !IsBinaryOnPath("brew") {

		fmt.Printf("⚙️  Homebrew is not installed, attempting to install...\n")

		if err := run("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"); err != nil {
			return fmt.Errorf("❌ Failed to install Homebrew: %w", err)
		}

		fmt.Printf("✅ Homebrew has been installed successfully\n")
	} else {
		fmt.Printf("✅ Homebrew is already installed\n")
	}

	return run("brew", "install", packageName)
}
