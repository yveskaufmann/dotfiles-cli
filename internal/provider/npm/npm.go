package npm

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
	"yv35.com/dotfiles-cli/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "npm"
}

func (p *Provider) Priority() int {
	return 100 // Default priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.NPM) > 0
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no npm packages defined
	if len(group.NPM) == 0 {
		return nil
	}

	// Check if npm is available
	if !sh.IsBinaryOnPath("npm") {
		onComplete(types.TaskResult{
			Name:   "npm-availability-check",
			Status: types.StatusSkipped,
			Error:  fmt.Errorf("npm is not installed or not found in PATH"),
		})
		return nil
	}

	// Batch install with verification
	return p.batchInstallWithVerification(group.NPM, onComplete)
}

func (p *Provider) batchInstallWithVerification(packages []config.NPMSpec, onComplete types.OnTaskComplete) error {
	// Filter out already installed packages
	var toInstall []string
	for _, npmPkg := range packages {
		pkg := npmPkg.Name

		// Check if already installed globally
		if err := sh.RunShell(fmt.Sprintf("npm list -g %s > /dev/null 2>&1", pkg)); err == nil {
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
	fmt.Printf("⚙️  Installing %d npm packages: %s\n", len(toInstall), strings.Join(toInstall, ", "))

	// Batch install
	args := append([]string{"install", "-g"}, toInstall...)
	if err := sh.Run("npm", args...); err != nil {
		// On batch install failure, report all as failed
		for _, pkg := range toInstall {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("batch install failed: %w", err),
			})
		}
		return fmt.Errorf("failed to batch install npm packages: %w", err)
	}

	// Post-verification: verify each package individually
	for _, pkg := range toInstall {
		if err := sh.RunShell(fmt.Sprintf("npm list -g %s > /dev/null 2>&1", pkg)); err != nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("verification failed: %w", err),
			})
			return fmt.Errorf("failed to verify npm package %s: %w", pkg, err)
		}

		onComplete(types.TaskResult{
			Name:   pkg,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
