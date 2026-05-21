package apt

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles-cli/internal/types"
	"yv35.com/dotfiles-cli/internal/util/sh"
)

func (p *Provider) batchInstallWithVerification(packages []string, onComplete types.OnTaskComplete) error {
	if len(packages) == 0 {
		return nil
	}

	if err := p.ensureAptUpdate(false); err != nil {
		return err
	}

	// Show progress message
	fmt.Printf("⚙️  Installing %d apt packages: %s\n", len(packages), strings.Join(packages, ", "))

	// Execute batch install
	args := append([]string{"apt", "install", "-y"}, packages...)
	err := sh.Run("sudo", args...)

	// Post-verification: Check each package individually
	for _, pkg := range packages {
		if p.isAptInstalled(pkg) {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSuccess,
				Error:  nil,
			})
		} else {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("package not found after installation"),
			})
		}
	}

	// Return error only if the batch command itself failed critically
	if err != nil {
		return fmt.Errorf("batch apt install encountered errors: %w", err)
	}

	return nil
}
