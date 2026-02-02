package types

import "yv35.com/dotfiles/internal/config"

type Provider interface {
	// Identifier returns the unique identifier of the provider (e.g., "github", "gitlab").
	ID() string

	// Priority returns the execution priority of the provider.
	// Lower numbers execute first. Typical values:
	// - 10: System package managers (apt, brew)
	// - 50: Version managers (nvm, sdkman)
	// - 100: Application installers (default)
	Priority() int

	// Install the installation of dependency group and returns a task Result.
	Install(config config.DependencyGroup, onComplete OnTaskComplete) error

	// HasConfig checks if the given dependency group contains configuration for this provider.
	HasConfig(group config.DependencyGroup) bool
}

type Setupable interface {
	// Setup performs any necessary setup for the provider before installations.
	// It will be invoked once by orchestrator before attempting any installations.
	Setup() error
}

type TearDownable interface {
	// TearDown performs any necessary cleanup after all installations are done.
	// It will be invoked once by orchestrator after all installations are complete.
	TearDown() error
}
