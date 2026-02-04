package apt

import (
	"fmt"

	"yv35.com/dotfiles/internal/util/sh"
)

func (p *Provider) ensureAptUpdate(force bool) error {
	if !p.aptUpdated || force {
		if err := sh.Run("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt index: %w", err)
		}
		p.aptUpdated = true
	}
	return nil
}
