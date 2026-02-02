package pipx

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/provider"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "pipx"
}

func (p *Provider) Priority() int {
	return 100 // Default priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.Pipx) > 0
}

func (p *Provider) Setup() error {
	if sh.IsBinaryOnPath("pipx") {
		return nil
	}

	// pipx not installed, attempt to install via system package manager
	if err := provider.InstallSystemPackage(provider.SystemPackageNames{
		Name: "pipx",
		Apt:  "pipx",
		Brew: "pipx",
	}); err != nil {
		return fmt.Errorf("failed to install pipx: %w", err)
	}

	if err := sh.RunShell("pipx ensurepath"); err != nil {
		return fmt.Errorf("failed to ensure pipx path: %w", err)
	}

	return nil
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no pipx packages defined
	if len(group.Pipx) == 0 {
		return nil
	}

	for _, pipxPkg := range group.Pipx {
		pkg := pipxPkg.Name

		// Check if already installed
		err := sh.RunShell("pipx list --short 2>/dev/null | awk '{print $1 }' | grep -xq -- \"" + pkg + "\"")
		if err == nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}

		// Install the package
		err = sh.Run("pipx", "install", pkg)
		if err != nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  err,
			})
			return fmt.Errorf("failed to install pipx package %s: %w", pkg, err)
		}

		onComplete(types.TaskResult{
			Name:   pkg,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
