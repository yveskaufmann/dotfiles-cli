package provider

import (
	"fmt"
	"runtime"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
)

type SystemPackageNames struct {
	Name string
	Apt  string
	Brew string
}

func InstallSystemPackage(pkgs SystemPackageNames) error {

	providerId := ""
	switch runtime.GOOS {
	case "linux":
		providerId = "apt"
	case "darwin":
		providerId = "brew"
	default:
		return fmt.Errorf("unsupported OS for installing %s with system package manager", pkgs.Name)
	}

	registry := NewRegistry()
	provider, exists := registry.GetProvider(providerId)
	if !exists {
		return fmt.Errorf("provider %s not found", providerId)
	}

	if err := SetupProvider(provider); err != nil {
		return err
	}

	dependency := config.DependencyGroup{
		Name: pkgs.Name,
		Apt:  []string{pkgs.Apt},
		Brew: []config.BrewSpec{
			{
				Name: pkgs.Brew,
			},
		},
	}

	OnComplete := func(result types.TaskResult) {
		if result.Error != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", dependency.Name, result.Error)
		} else {
			fmt.Printf("✅ %s installed successfully\n", dependency.Name)
		}
	}

	return provider.Install(dependency, OnComplete)
}
