package tools

import (
	"fmt"

	_os "yv35.com/dotfiles/internal/os"
)

func InstallSnapPackages(packages []string) error {
	if !_os.IsUbuntu() {
		fmt.Printf("⚠️  Snap packages installation is skipped on non Ubuntu systems.\n")
		return nil
	}

	for _, pkg := range packages {
		if err := RunShell("snap list " + pkg + " &> /dev/null"); err == nil {
			fmt.Printf("✅ Snap package %s is already installed\n", pkg)
			continue
		}

		fmt.Printf("⚙️  Installing snap package: %s\n", pkg)
		err := Run("sudo", "snap", "install", pkg)
		if err != nil {
			return fmt.Errorf("❌ Failed to install snap package %s: %w", pkg, err)
		}
		fmt.Printf("✅ Snap package %s installed successfully\n", pkg)
	}

	return nil
}
