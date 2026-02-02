package apt

import (
	"fmt"

	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
	"yv35.com/dotfiles/internal/util/stringutils"
)

func (p *Provider) isURL(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}

func (p *Provider) installDebFromURL(url string, onComplete types.OnTaskComplete) error {

	url = stringutils.ResolvePlaceholders(url)

	err := sh.RunShell(fmt.Sprintf("curl -L -o /tmp/package.deb %s && sudo apt install -y /tmp/package.deb && rm /tmp/package.deb", url))

	if err != nil {
		onComplete(types.TaskResult{
			Name:   url,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to install deb from url: %w", err)
	}

	if err := p.ensureAptUpdate(true); err != nil {
		onComplete(types.TaskResult{
			Name:   url,
			Status: types.StatusFailed,
			Error:  err,
		})
		return err
	}

	onComplete(types.TaskResult{
		Name:   url,
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return nil
}
