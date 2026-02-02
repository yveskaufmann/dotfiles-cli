package provider

import (
	"fmt"

	"yv35.com/dotfiles/internal/types"
)

func SetupProvider(provider types.Provider) error {
	if setupable, ok := provider.(types.Setupable); ok {
		if err := setupable.Setup(); err != nil {
			return fmt.Errorf("failed to setup provider %s: %w", provider.ID(), err)
		}
	}
	return nil
}
