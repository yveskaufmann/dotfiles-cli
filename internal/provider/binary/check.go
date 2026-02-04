package binary

import (
	"fmt"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
)

// isAlreadyInstalled checks if a binary spec is already installed
func (p *Provider) isAlreadyInstalled(spec config.BinarySpec) bool {
	if len(spec.Binaries) > 0 {
		for _, b := range spec.Binaries {
			if err := sh.RunShell("type " + b + " > /dev/null 2>&1"); err != nil {
				return false
			}
		}
		return true
	}
	return sh.RunShell("type "+spec.Name+" > /dev/null 2>&1") == nil
}

// filterInstalled separates already installed binaries from those that need installation
func (p *Provider) filterInstalled(specs []config.BinarySpec, onComplete types.OnTaskComplete) []config.BinarySpec {
	var toInstall []config.BinarySpec
	for _, spec := range specs {
		if p.isAlreadyInstalled(spec) {
			onComplete(types.TaskResult{
				Name:   spec.Name,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
		} else {
			toInstall = append(toInstall, spec)
		}
	}
	return toInstall
}
