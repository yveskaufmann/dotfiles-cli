package tool

import (
	"fmt"

	"yv35.com/dotfiles/internal/util/osutil"
	"yv35.com/dotfiles/internal/util/sh"
)

func InstallSnapPackages(packages []string) error {

	if !osutil.IsUbuntu() {
		fmt.Printf("⚠️  Snap packages installation is skipped on non Ubuntu systems.\n")
		return nil
	}

	if !sh.IsBinaryOnPath("snap") {
		fmt.Printf("⚠️ snap is not installed or not found in PATH . Skipping snap packages installation.\n")
		return nil
	}

	for _, pkg := range packages {
		if err := sh.RunShell("snap list " + pkg + " &> /dev/null"); err == nil {
			fmt.Printf("✅ Snap package %s is already installed\n", pkg)
			continue
		}

		fmt.Printf("⚙️  Installing snap package: %s\n", pkg)
		err := sh.RunShell("sudo snap install " + pkg)
		if err != nil {
			return fmt.Errorf("❌ Failed to install snap package %s: %w", pkg, err)
		}
		fmt.Printf("✅ Snap package %s installed successfully\n", pkg)
	}

	return nil
}
