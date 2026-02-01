package tool

import (
	"fmt"

	"yv35.com/dotfiles/internal/util/sh"
)

func InstallNPMPackages(packages []string) error {
	for _, pkg := range packages {
		if err := InstallNPM(pkg); err != nil {
			return err
		}
	}
	return nil
}

func InstallNPM(pkg string) error {
	if !sh.IsBinaryOnPath("npm") {
		return fmt.Errorf("npm is not installed, cannot install package %s", pkg)
	}

	// check if already installed globally
	if err := sh.RunShell(fmt.Sprintf("npm list -g %s > /dev/null 2>&1", pkg)); err == nil {
		fmt.Printf("✅ NPM package %s is already installed\n", pkg)
		return nil
	}

	fmt.Printf("⚙️  Installing NPM package: %s\n", pkg)
	if err := sh.Run("npm", "install", "-g", pkg); err != nil {
		return fmt.Errorf("failed to install npm package %s: %w", pkg, err)
	}
	fmt.Printf("✅ NPM package %s installed successfully\n", pkg)
	return nil
}
