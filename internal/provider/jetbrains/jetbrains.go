package jetbrains

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
)

// Provider implements types.Provider for JetBrains IDEs.  It downloads IDEs
// from the JetBrains public API and installs them under /opt/<ide-dir>.
// Existing installations are replaced when the installed version differs from
// the desired version (pinned or latest).
type Provider struct{}

// NewProvider creates and returns a new JetBrains Provider.
func NewProvider() *Provider {
	return &Provider{}
}

// ID returns the unique identifier used in the provider registry.
func (p *Provider) ID() string {
	return "jetbrains"
}

// Priority returns the execution priority.  100 = application installer tier.
func (p *Provider) Priority() int {
	return 100
}

// HasConfig reports whether the dependency group contains at least one
// JetBrains IDE spec.
func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.Jetbrains) > 0
}

// Install checks each IDE spec, skips those that are already at the desired
// version, and installs (or upgrades) the rest concurrently.
func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	if len(group.Jetbrains) == 0 {
		return nil
	}

	fmt.Printf("🔍 Checking installed versions for %d JetBrains IDE(s)...\n", len(group.Jetbrains))

	// Filter out already up-to-date IDEs; collect install candidates.
	toInstall := p.filterInstalled(group.Jetbrains, onComplete)
	if len(toInstall) == 0 {
		return nil
	}

	names := make([]string, len(toInstall))
	for i, c := range toInstall {
		if c.Spec.Name != "" {
			names[i] = fmt.Sprintf("%s (%s)", c.Spec.Name, c.Release.Version)
		} else {
			names[i] = fmt.Sprintf("%s (%s)", c.Spec.IDE, c.Release.Version)
		}
	}
	fmt.Printf(
		"⚙️  Installing %d JetBrains IDE(s) concurrently (max %d parallel): %s\n",
		len(toInstall), maxConcurrentDownloads, strings.Join(names, ", "),
	)

	return p.processWithWorkerPool(toInstall, onComplete)
}
