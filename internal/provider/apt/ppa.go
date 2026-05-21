package apt

import (
	"fmt"
	"log/slog"
	"os"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
)

// addPPASource adds the PPA source without installing packages
// Returns: (sourceWasAdded, packagesToInstall, error)
func (p *Provider) addPPASource(ppa config.PPASpec, onComplete types.OnTaskComplete) (bool, []string, error) {
	name := ppa.Name

	// Check if source already exists
	sourcesPath := fmt.Sprintf("/etc/apt/sources.list.d/%s.sources", ppa.SourceName)
	if _, err := os.Stat(sourcesPath); err == nil {
		// Source already exists, just return packages
		return false, ppa.Pkgs, nil
	}

	// Handle GPG Key
	keyringPath, err := p.installKeyFile(ppa, onComplete)
	if err != nil {
		return false, nil, err
	}

	err, sourceData := ReadPPASourceFromPPASpec(ppa, keyringPath)
	if err != nil {
		onComplete(types.TaskResult{
			Name:   "ppa:" + name + " source data",
			Status: types.StatusFailed,
			Error:  err,
		})
		return false, nil, err
	}

	sourceFileContent := RenderPPASourceFile(sourceData)

	slog.Info(sourceFileContent)

	if err := WritePPASourceFile(sourcesPath, ppa, sourceFileContent); err != nil {
		onComplete(types.TaskResult{
			Name:   "ppa:" + name + " source file",
			Status: types.StatusFailed,
			Error:  err,
		})
		return false, nil, err
	}

	onComplete(types.TaskResult{
		Name:   "ppa:" + name + " source file",
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return true, ppa.Pkgs, nil
}
