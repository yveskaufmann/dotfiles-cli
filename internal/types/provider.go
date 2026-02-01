package types

import "yv35.com/dotfiles/internal/config"

type Provider interface {
	// Identifier returns the unique identifier of the provider (e.g., "github", "gitlab").
	ID() string

	// Install the installation of dependency group and returns a task Result.
	Install(config config.DependencyGroup, onComplete OnTaskComplete) error
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
