package executor

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/provider"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/osutil"
)

type ToolInstallExecutor struct {
	Groups           []config.DependencyGroup
	Providers        *provider.Registry
	EnabledProviders map[string]bool
}

func NewToolInstallExecutor(config *config.Config) *ToolInstallExecutor {
	return &ToolInstallExecutor{
		Groups:           config.Groups,
		Providers:        provider.NewRegistry(),
		EnabledProviders: nil, // nil means all enabled
	}
}

func (e *ToolInstallExecutor) SetEnabledProviders(providers []string) {
	if len(providers) == 0 {
		e.EnabledProviders = nil
		return
	}

	e.EnabledProviders = make(map[string]bool)
	for _, p := range providers {
		e.EnabledProviders[p] = true
	}
}

func (e *ToolInstallExecutor) isProviderEnabled(providerID string) bool {
	if e.EnabledProviders == nil {
		return true
	}
	return e.EnabledProviders[providerID]
}

func (e *ToolInstallExecutor) Setup() error {
	for _, providerInstance := range e.Providers.List() {
		if !e.isProviderEnabled(providerInstance.ID()) {
			continue
		}

		if setupable, ok := providerInstance.(types.Setupable); ok {
			if err := setupable.Setup(); err != nil {
				return fmt.Errorf("failed to setup provider %s: %w", providerInstance.ID(), err)
			}
		}
	}

	return nil
}

func (e *ToolInstallExecutor) Execute() error {
	for _, group := range e.Groups {
		if group.Systems != "" && !osutil.Is(osutil.OSType(group.Systems)) {
			fmt.Printf("⏭️  Skipping group %s (system mismatch: %s)\n", group.Name, group.Systems)
			continue
		}

		fmt.Printf("🚀 Processing group: %s\n", group.Name)
		if err := e.executeGroup(group); err != nil {
			return fmt.Errorf("failed to execute group %s: %w", group.Name, err)
		}
	}
	return nil
}

func (e *ToolInstallExecutor) executeGroup(group config.DependencyGroup) error {

	for _, providerInstance := range e.Providers.List() {
		if !e.isProviderEnabled(providerInstance.ID()) {
			continue
		}

		onProgress := func(result types.TaskResult) {
			switch result.Status {
			case types.StatusSuccess:
				fmt.Printf("✅ [%s] %s installed successfully\n", providerInstance.ID(), result.Name)
			case types.StatusSkipped:
				fmt.Printf("⏭️  [%s] %s installation skipped: %v\n", providerInstance.ID(), result.Name, result.Error)
			case types.StatusFailed:
				fmt.Printf("❌ [%s] %s installation failed: %v\n", providerInstance.ID(), result.Name, result.Error)
			}
		}

		err := providerInstance.Install(group, onProgress)
		if err != nil {
			return fmt.Errorf("provider %s failed to install for group %s: %w", providerInstance.ID(), group.Name, err)
		}
	}

	return nil
}
