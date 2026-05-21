package apt

import (
	"fmt"
	"log/slog"
	"os"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
	"yv35.com/dotfiles-cli/internal/util/sh"
)

func (p *Provider) ensureGPGInitialized() error {

	if !p.isAptInstalled("gnupg") {
		sh.RunShell("sudo apt-get install gnupg -yy")
	}

	pubringKeyPath := os.ExpandEnv("$HOME/.gnupg/pubring.kbx")
	if _, err := os.Stat(pubringKeyPath); os.IsNotExist(err) {
		if err := sh.RunShell("sudo gpg --list-keys > /dev/null 2>&1"); err != nil {
			return fmt.Errorf("failed to initialize GPG: %w", err)
		}
		slog.Info("Initialized GPG keyring at", "path", pubringKeyPath)
	}
	return nil
}

func (p *Provider) installKeyFile(ppa config.PPASpec, onComplete types.OnTaskComplete) (string, error) {

	keyringPath := ""

	if ppa.Key == "" && (ppa.KeyServer == "" && ppa.KeyID == "") {
		return keyringPath, nil
	}

	if err := p.ensureGPGInitialized(); err != nil {
		slog.Error("Failed to initialize GPG", "error", err)
		return "", err
	}

	keyringPath = fmt.Sprintf("/usr/share/keyrings/%s-keyring.gpg", ppa.SourceName)

	_, err := os.Stat(keyringPath)
	if err == nil {
		onComplete(types.TaskResult{
			Name:   "keyring - ppa:" + ppa.Name,
			Status: types.StatusSkipped,
			Error:  fmt.Errorf("keyfile already installed: %s", keyringPath),
		})
		return keyringPath, nil
	}

	if !os.IsNotExist(err) {
		onComplete(types.TaskResult{
			Name:   "keyring - ppa:" + ppa.Name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return "", fmt.Errorf("failed to stat keyring file: %w", err)
	}

	switch {
	case ppa.Key != "":
		if err := p.downloadKeyFile(keyringPath, ppa, onComplete); err != nil {
			return "", err
		}
	case ppa.KeyServer != "" && ppa.KeyID != "":
		if err := sh.RunShell(fmt.Sprintf("sudo gpg --no-default-keyring --keyring %s --keyserver %s --recv-keys %s", keyringPath, ppa.KeyServer, ppa.KeyID)); err != nil {
			onComplete(types.TaskResult{
				Name:   "ppa:" + ppa.Name,
				Status: types.StatusFailed,
				Error:  err,
			})
			return "", fmt.Errorf("failed to download GPG key from server: %w", err)
		}
	default:
		return "", fmt.Errorf("no key information provided for PPA: %s", ppa.Name)
	}

	onComplete(types.TaskResult{
		Name:   "keyring - ppa:" + ppa.Name,
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return keyringPath, nil
}

func (p *Provider) downloadKeyFile(keyringPath string, ppa config.PPASpec, onComplete types.OnTaskComplete) error {
	if ppa.Key == "" {
		return fmt.Errorf("missing key file to downlad keyring file")
	}

	if err := sh.RunShell(fmt.Sprintf("curl -sSL %s | sudo gpg --dearmor --yes --output %s", ppa.Key, keyringPath)); err != nil {
		onComplete(types.TaskResult{
			Name:   "ppa:" + ppa.Name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to download and dearmor GPG key: %w", err)
	}
	return nil
}
