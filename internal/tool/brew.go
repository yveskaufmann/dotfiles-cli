package tool

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/util/osutil"
	"yv35.com/dotfiles/internal/util/sh"
)

func InstallBrewPackages(dependency config.DependencyGroup) error {

	if len(dependency.Brew) == 0 {
		return nil
	}

	if !osutil.IsMac() {
		fmt.Printf("⚠️  Brew packages installation is skipped on non Mac systems.\n")
		return nil
	}

	if !sh.IsBinaryOnPath("brew") {
		fmt.Printf("⚙️  Homebrew is not installed, attempting to install...\n")

		if err := sh.Run("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"); err != nil {
			return fmt.Errorf("❌ Failed to install Homebrew: %w.  Skipping brew packages installation.", err)
		}
		fmt.Printf("✅ Homebrew has been installed successfully\n")
	}

	for _, formula := range dependency.Brew {
		pkg := formula.Name

		brewArgs := ""
		if formula.Cask {
			brewArgs += "--cask "
		}

		if err := sh.RunShell(fmt.Sprintf("brew list %s %s &> /dev/null", brewArgs, pkg)); err == nil {
			fmt.Printf("✅ Brew package %s is already installed\n", pkg)
			continue
		}

		fmt.Printf("⚙️  Installing brew package: %s\n", pkg)
		err := sh.RunShell(fmt.Sprintf("HOMEBREW_NO_AUTO_UPDATE=1 brew install %s %s", brewArgs, pkg))
		if err != nil {
			return fmt.Errorf("❌ Failed to install brew package %s: %w", pkg, err)
		}
		fmt.Printf("✅ Brew package %s installed successfully\n", pkg)
	}

	return nil
}

func InstallBrewTap(dependency config.DependencyGroup) error {

	if len(dependency.BrewTapSpec) == 0 {
		return nil
	}

	return nil
}

func Install_(packageName string) error {
	if !sh.IsBinaryOnPath("brew") {

		fmt.Printf("⚙️  Homebrew is not installed, attempting to install...\n")

		if err := sh.Run("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"); err != nil {
			return fmt.Errorf("❌ Failed to install Homebrew: %w", err)
		}

		fmt.Printf("✅ Homebrew has been installed successfully\n")
	} else {
		fmt.Printf("✅ Homebrew is already installed\n")
	}

	return sh.Run("brew", "install", packageName)
}
