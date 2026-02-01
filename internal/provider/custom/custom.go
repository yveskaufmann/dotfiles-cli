package custom

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "custom"
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no custom packages defined
	if len(group.Custom) == 0 {
		return nil
	}

	for _, customSpec := range group.Custom {
		name := customSpec.Name

		// Check if already installed
		if customSpec.InstallCheck != "" {
			if err := sh.RunShell(customSpec.InstallCheck); err == nil {
				onComplete(types.TaskResult{
					Name:   name,
					Status: types.StatusSkipped,
					Error:  fmt.Errorf("install check passed"),
				})
				continue
			}
		}

		// Install
		err := sh.RunShell(customSpec.Install)
		if err != nil {
			onComplete(types.TaskResult{
				Name:   name,
				Status: types.StatusFailed,
				Error:  err,
			})
			return fmt.Errorf("failed to install %s: %w", name, err)
		}

		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
