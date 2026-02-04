package sdkman

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "sdkman"
}

func (p *Provider) Priority() int {
	return 50 // Version manager - medium priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.Sdkman) > 0
}

func (p *Provider) ensureDependencies() error {
	// Check if zip and unzip are installed
	_, zipErr := sh.RunShellOutput("which zip 2>/dev/null")
	_, unzipErr := sh.RunShellOutput("which unzip 2>/dev/null")

	if zipErr == nil && unzipErr == nil {
		return nil // Both already installed
	}

	fmt.Println("⚙️  Installing SdkMan dependencies (zip, unzip)...")

	switch runtime.GOOS {
	case "linux":
		// Use apt for Ubuntu/Debian
		if err := sh.RunShell("sudo apt-get update -qq && sudo apt-get install -y zip unzip"); err != nil {
			return fmt.Errorf("failed to install zip/unzip via apt: %w", err)
		}
		fmt.Println("✓ SdkMan dependencies installed")
	case "darwin":
		// macOS usually has these preinstalled, but check brew if needed
		if zipErr != nil || unzipErr != nil {
			// Try brew but don't fail if it doesn't work (likely preinstalled)
			sh.RunShell("brew install zip unzip 2>/dev/null || true")
		}
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return nil
}

func (p *Provider) Setup() error {
	// Ensure zip and unzip are installed
	if err := p.ensureDependencies(); err != nil {
		return err
	}

	sdkmanDir := p.getSdkmanDir()

	// Install SdkMan if not present
	if !p.isSdkmanInstalled(sdkmanDir) {
		fmt.Println("⚙️  Installing SdkMan...")

		// Download and run install script
		installCmd := `curl -s "https://get.sdkman.io" | bash`

		if err := sh.RunShell(installCmd); err != nil {
			return fmt.Errorf("failed to install sdkman: %w", err)
		}

		// Run self-update
		updateCmd := p.sdkCmd("sdk selfupdate")
		if err := sh.RunShell(updateCmd); err != nil {
			fmt.Printf("⚠️  Warning: Failed to self-update sdkman: %v\n", err)
		}

		fmt.Println("✓ SdkMan installed successfully")
	} else {
		fmt.Println("✓ SdkMan is already installed")

		// Update existing installation
		updateCmd := p.sdkCmd("sdk selfupdate")
		if err := sh.RunShell(updateCmd); err != nil {
			fmt.Printf("⚠️  Warning: Failed to self-update sdkman: %v\n", err)
		}
	}

	return nil
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no sdkman configurations defined
	if len(group.Sdkman) == 0 {
		return nil
	}

	for _, sdkmanSpec := range group.Sdkman {
		candidate := sdkmanSpec.Candidate
		defaultVersion, versions := p.resolveVersions(sdkmanSpec)

		// Install each version
		for _, version := range versions {
			if err := p.installCandidate(candidate, version, onComplete); err != nil {
				return err
			}
		}

		// Set the default version
		if err := p.setDefaultVersion(candidate, defaultVersion, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) getSdkmanDir() string {
	if sdkmanDir := os.Getenv("SDKMAN_DIR"); sdkmanDir != "" {
		return sdkmanDir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".sdkman")
}

func (p *Provider) isSdkmanInstalled(sdkmanDir string) bool {
	sdkmanScript := filepath.Join(sdkmanDir, "bin", "sdkman-init.sh")
	_, err := os.Stat(sdkmanScript)
	return err == nil
}

func (p *Provider) sdkCmd(command string) string {
	sdkmanDir := p.getSdkmanDir()
	return fmt.Sprintf(`export SDKMAN_DIR="%s" && [ -s "$SDKMAN_DIR/bin/sdkman-init.sh" ] && source "$SDKMAN_DIR/bin/sdkman-init.sh" && %s`, sdkmanDir, command)
}

// resolveVersions determines the default version and all versions to install
// Logic:
// 1. If Version is set, it's the default
// 2. Otherwise, first item in Versions array is the default
// 3. All versions (including default) are returned in the versions slice
func (p *Provider) resolveVersions(spec config.SdkmanSpec) (defaultVersion string, versions []string) {
	versionSet := make(map[string]bool)
	result := []string{}

	// Determine default version
	if spec.Version != "" {
		defaultVersion = spec.Version
		versionSet[defaultVersion] = true
		result = append(result, defaultVersion)
	} else if len(spec.Versions) > 0 {
		defaultVersion = spec.Versions[0]
		versionSet[defaultVersion] = true
		result = append(result, defaultVersion)
	}

	// Add additional versions from the Versions array
	for _, v := range spec.Versions {
		if !versionSet[v] {
			versionSet[v] = true
			result = append(result, v)
		}
	}

	// If we have a Version field AND Versions array, add the Version field versions too
	if spec.Version != "" {
		for _, v := range spec.Versions {
			if !versionSet[v] {
				versionSet[v] = true
				result = append(result, v)
			}
		}
	}

	return defaultVersion, result
}

func (p *Provider) isCandidateInstalled(candidate, version string) bool {
	sdkmanDir := p.getSdkmanDir()
	versionPath := filepath.Join(sdkmanDir, "candidates", candidate, version)
	_, err := os.Stat(versionPath)
	return err == nil
}

func (p *Provider) installCandidate(candidate, version string, onComplete types.OnTaskComplete) error {
	taskName := fmt.Sprintf("%s@%s", candidate, version)

	// Check if already installed
	if p.isCandidateInstalled(candidate, version) {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusSkipped,
		})
		return nil
	}

	// Install the version (< /dev/null to avoid interactive prompts)
	installCmd := p.sdkCmd(fmt.Sprintf("sdk install %s %s < /dev/null", candidate, version))

	if err := sh.RunShell(installCmd); err != nil {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusFailed,
			Error:  fmt.Errorf("failed to install %s version %s: %w", candidate, version, err),
		})
		return fmt.Errorf("failed to install %s version %s: %w", candidate, version, err)
	}

	onComplete(types.TaskResult{
		Name:   taskName,
		Status: types.StatusSuccess,
	})

	return nil
}

func (p *Provider) setDefaultVersion(candidate, version string, onComplete types.OnTaskComplete) error {
	taskName := fmt.Sprintf("%s@%s (set as default)", candidate, version)

	// Get current default
	currentDefault, _ := p.getCurrentDefault(candidate)

	// Skip if already set
	if strings.TrimSpace(currentDefault) == version {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusSkipped,
		})
		return nil
	}

	// Set default
	setDefaultCmd := p.sdkCmd(fmt.Sprintf("sdk default %s %s", candidate, version))

	if err := sh.RunShell(setDefaultCmd); err != nil {
		onComplete(types.TaskResult{
			Name:   taskName,
			Status: types.StatusFailed,
			Error:  fmt.Errorf("failed to set default %s version to %s: %w", candidate, version, err),
		})
		return fmt.Errorf("failed to set default %s version to %s: %w", candidate, version, err)
	}

	onComplete(types.TaskResult{
		Name:   taskName,
		Status: types.StatusSuccess,
	})

	return nil
}

func (p *Provider) getCurrentDefault(candidate string) (string, error) {
	cmd := p.sdkCmd(fmt.Sprintf("sdk current %s 2>/dev/null | awk '{print $NF}'", candidate))

	output, err := sh.RunShellOutput(cmd)
	return strings.TrimSpace(output), err
}
