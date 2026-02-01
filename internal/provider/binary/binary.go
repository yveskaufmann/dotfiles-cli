package binary

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "binary"
}

func (p *Provider) Priority() int {
	return 100 // Default priority
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no binary packages defined
	if len(group.Binary) == 0 {
		return nil
	}

	// Filter out already installed binaries
	toInstall := p.filterInstalled(group.Binary, onComplete)

	if len(toInstall) == 0 {
		return nil
	}

	// Show progress message
	names := make([]string, len(toInstall))
	for i, spec := range toInstall {
		names[i] = spec.Name
	}
	fmt.Printf("⚙️  Installing %d binaries concurrently (max %d parallel): %s\n",
		len(toInstall), maxConcurrentDownloads, strings.Join(names, ", "))

	// Process with worker pool
	return p.processWithWorkerPool(toInstall, onComplete)
}
