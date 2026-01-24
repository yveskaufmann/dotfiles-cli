package tools

import (
	"fmt"
)

func EnsurePipx() error {
	if IsBinaryOnPath("pipx") {
		return nil
	}

	fmt.Printf("⚙️ pipx is not installed, attempting to install...\n")

	if err := Install("pipx"); err != nil {
		return fmt.Errorf("failed to install pipx: %w", err)
	}

	if err := RunShell("pipx ensurepath"); err != nil {
		return fmt.Errorf("failed to ensure pipx path: %w", err)
	}

	fmt.Printf("✅ pipx has been installed successfully\n")

	return nil
}

func InstallPipx(pkg string) error {
	if err := EnsurePipx(); err != nil {
		return err
	}

	err := RunShell("pipx list --short 2>/dev/null | awk '{print $1 }' | grep -xq -- \"" + pkg + "\"")
	if err == nil {
		fmt.Printf("✅ pipx package %s is already installed\n", pkg)
		return nil
	}

	fmt.Printf("⚙️  Installing pipx package: %s\n", pkg)
	err = Run("pipx", "install", pkg)
	if err != nil {
		return fmt.Errorf("❌ Failed to install pipx package %s: %w", pkg, err)
	}
	fmt.Printf("✅ Pipx package %s installed successfully\n", pkg)

	return nil
}
