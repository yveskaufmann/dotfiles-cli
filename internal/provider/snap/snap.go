package snap

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
	"yv35.com/dotfiles-cli/internal/util/osutil"
	"yv35.com/dotfiles-cli/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "snap"
}

func (p *Provider) Priority() int {
	return 100 // Default priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.Snap) > 0
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no snap packages defined
	if len(group.Snap) == 0 {
		return nil
	}

	// Skip if not Ubuntu
	if !osutil.IsUbuntu() {
		return nil
	}

	// Check if snap is available
	if !sh.IsBinaryOnPath("snap") {
		onComplete(types.TaskResult{
			Name:   "snap-availability-check",
			Status: types.StatusSkipped,
			Error:  fmt.Errorf("snap is not installed or not found in PATH"),
		})
		return nil
	}

	// Batch install with verification
	return p.batchInstallWithVerification(group.Snap, onComplete)
}

func (p *Provider) batchInstallWithVerification(packages []config.SnapSpec, onComplete types.OnTaskComplete) error {
	// Filter out already installed packages
	var toInstall []string
	for _, snapPkg := range packages {
		pkg := snapPkg.Name

		// Check if already installed
		if err := sh.RunShell("snap list " + pkg + " &> /dev/null"); err == nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
		} else {
			toInstall = append(toInstall, pkg)
		}
	}

	// Nothing to install
	if len(toInstall) == 0 {
		return nil
	}

	// Show progress message
	fmt.Printf("⚙️  Installing %d snap packages: %s\n", len(toInstall), strings.Join(toInstall, ", "))

	// Batch install
	installCmd := "sudo snap install " + strings.Join(toInstall, " ")
	if err := sh.RunShell(installCmd); err != nil {
		// On batch install failure, report all as failed
		for _, pkg := range toInstall {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("batch install failed: %w", err),
			})
		}
		return fmt.Errorf("failed to batch install snap packages: %w", err)
	}

	// Post-verification: verify each package individually
	for _, pkg := range toInstall {
		if err := sh.RunShell("snap list " + pkg + " &> /dev/null"); err != nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("verification failed: %w", err),
			})
			return fmt.Errorf("failed to verify snap package %s: %w", pkg, err)
		}

		onComplete(types.TaskResult{
			Name:   pkg,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
