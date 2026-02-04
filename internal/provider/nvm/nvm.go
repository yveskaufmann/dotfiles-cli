package nvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/provider/github"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "nvm"
}

func (p *Provider) Priority() int {
	return 50 // Version manager - medium priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.NVM) > 0
}

func (p *Provider) Setup() error {
	nvmDir := p.getNVMDir()

	// Install NVM if not present
	if !p.isNVMInstalled(nvmDir) {
		fmt.Println("⚙️  Installing NVM...")

		// Get latest NVM release tag
		tag, err := github.GetLatestRelease("nvm-sh/nvm")
		if err != nil {
			return fmt.Errorf("failed to get latest nvm release: %w", err)
		}

		// Download and run install script
		installURL := fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/%s/install.sh", tag)
		installCmd := fmt.Sprintf("curl -o- %s | bash", installURL)

		if err := sh.RunShell(installCmd); err != nil {
			return fmt.Errorf("failed to install nvm: %w", err)
		}

		fmt.Printf("✓ NVM %s installed successfully\n", tag)
	} else {
		fmt.Println("✓ NVM is already installed")
	}

	return nil
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no nvm configurations defined
	if len(group.NVM) == 0 {
		return nil
	}

	nvmDir := p.getNVMDir()

	for _, nvmSpec := range group.NVM {
		// Resolve "lts/latest" to actual latest LTS version
		defaultVersion, err := p.resolveLtsVersion(nvmDir, nvmSpec.Default)
		if err != nil {
			return fmt.Errorf("failed to resolve default version %s: %w", nvmSpec.Default, err)
		}

		// Resolve versions list
		resolvedVersions := []string{}
		for _, v := range nvmSpec.Versions {
			resolved, err := p.resolveLtsVersion(nvmDir, v)
			if err != nil {
				return fmt.Errorf("failed to resolve version %s: %w", v, err)
			}
			resolvedVersions = append(resolvedVersions, resolved)
		}

		// Deduplicate versions: ensure default is in the list
		versions := p.deduplicateVersions(defaultVersion, resolvedVersions)

		// Install each version
		for _, version := range versions {
			if err := p.installNodeVersion(nvmDir, version, onComplete); err != nil {
				return err
			}
		}

		// Set the default version
		if err := p.setDefaultVersion(nvmDir, defaultVersion, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) getNVMDir() string {
	if nvmDir := os.Getenv("NVM_DIR"); nvmDir != "" {
		return nvmDir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".nvm")
}

func (p *Provider) isNVMInstalled(nvmDir string) bool {
	nvmScript := filepath.Join(nvmDir, "nvm.sh")
	_, err := os.Stat(nvmScript)
	return err == nil
}

func (p *Provider) resolveLtsVersion(nvmDir, version string) (string, error) {
	// If not "lts/latest", return as-is
	if version != "lts/latest" {
		return version, nil
	}

	// Ensure NVM is installed before trying to resolve
	if !p.isNVMInstalled(nvmDir) {
		return "", fmt.Errorf("NVM is not installed, cannot resolve lts/latest")
	}

	// Resolve "lts/latest" to the actual latest LTS version
	resolveCmd := fmt.Sprintf(`export NVM_DIR="%s" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && nvm version-remote --lts`, nvmDir)

	output, err := sh.RunShellOutput(resolveCmd)
	if err != nil {
		return "", fmt.Errorf("failed to resolve latest LTS version: %w", err)
	}

	resolved := strings.TrimSpace(output)
	if resolved == "" {
		return "", fmt.Errorf("could not determine latest LTS version")
	}

	return resolved, nil
}

func (p *Provider) deduplicateVersions(defaultVersion string, versions []string) []string {
	versionSet := make(map[string]bool)
	result := []string{}

	// Add default first
	versionSet[defaultVersion] = true
	result = append(result, defaultVersion)

	// Add other versions if not already present
	for _, v := range versions {
		if !versionSet[v] {
			versionSet[v] = true
			result = append(result, v)
		}
	}

	return result
}

func (p *Provider) isNodeVersionInstalled(nvmDir, version string) bool {
	// Source nvm and check if version is installed
	checkCmd := fmt.Sprintf(`
		export NVM_DIR="%s"
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
		nvm ls %s &>/dev/null
	`, nvmDir, version)

	err := sh.RunShell(checkCmd)
	return err == nil
}

func (p *Provider) installNodeVersion(nvmDir, version string, onComplete types.OnTaskComplete) error {
	taskName := fmt.Sprintf("node@%s", version)

	// Check if already installed
	if p.isNodeVersionInstalled(nvmDir, version) {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusSkipped,
		})
		return nil
	}

	// Install the version
	installCmd := fmt.Sprintf(`
		export NVM_DIR="%s"
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
		nvm install %s
	`, nvmDir, version)

	if err := sh.RunShell(installCmd); err != nil {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusFailed,
			Error:  fmt.Errorf("failed to install node version %s: %w", version, err),
		})
		return fmt.Errorf("failed to install node version %s: %w", version, err)
	}

	onComplete(types.TaskResult{
		Name:   taskName,
		Status: types.StatusSuccess,
	})

	return nil
}

func (p *Provider) setDefaultVersion(nvmDir, version string, onComplete types.OnTaskComplete) error {
	taskName := fmt.Sprintf("node@%s (set as default)", version)

	// Get current default
	currentDefault, _ := p.getCurrentDefault(nvmDir)

	// Skip if already set
	if strings.TrimSpace(currentDefault) == version {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusSkipped,
		})
		return nil
	}

	// Set default alias
	setDefaultCmd := fmt.Sprintf(`
		export NVM_DIR="%s"
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
		nvm alias default %s
	`, nvmDir, version)

	if err := sh.RunShell(setDefaultCmd); err != nil {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusFailed,
			Error:  fmt.Errorf("failed to set default node version to %s: %w", version, err),
		})
		return fmt.Errorf("failed to set default node version to %s: %w", version, err)
	}

	onComplete(types.TaskResult{
		Name:   taskName,
		Status: types.StatusSuccess,
	})

	return nil
}

func (p *Provider) getCurrentDefault(nvmDir string) (string, error) {
	cmd := fmt.Sprintf(`
		export NVM_DIR="%s"
		[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
		nvm alias default 2>/dev/null | awk '{print $NF}' | sed 's/[()]//g'
	`, nvmDir)

	output, err := sh.RunShellOutput(cmd)
	return strings.TrimSpace(output), err
}
