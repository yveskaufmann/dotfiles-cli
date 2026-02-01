package script

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
	return "script"
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no scripts defined
	if len(group.Script) == 0 {
		return nil
	}

	for _, scriptSpec := range group.Script {
		name := scriptSpec.Name
		script := scriptSpec.Script

		err := sh.RunScript(script)
		if err != nil {
			onComplete(types.TaskResult{
				Name:   name,
				Status: types.StatusFailed,
				Error:  err,
			})
			return fmt.Errorf("failed to run script %s: %w", name, err)
		}

		onComplete(types.TaskResult{
			Name:   name,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
