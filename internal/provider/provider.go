package provider

import (
	"fmt"
	"log/slog"

	"yv35.com/dotfiles/internal/types"
)

func SetupProvider(provider types.Provider) error {
	if setupable, ok := provider.(types.Setupable); ok {
		slog.Info("⚙️ Setting up provider", "provider", provider.ID())
		if err := setupable.Setup(); err != nil {
			return fmt.Errorf("failed to setup provider %s: %w", provider.ID(), err)
		}
		slog.Info("✅⚙️Finished setting up provider", "provider", provider.ID())
	}
	return nil
}
