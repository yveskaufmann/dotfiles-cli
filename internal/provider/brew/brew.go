package brew

import (
	"fmt"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/osutil"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return "brew"
}

func (p *Provider) Setup() error {
	// Skip if not macOS
	if !osutil.IsMac() {
		return nil
	}

	// Check if brew is already installed
	if sh.IsBinaryOnPath("brew") {
		return nil
	}

	// Install Homebrew
	if err := sh.Run("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}

	return nil
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if not macOS
	if !osutil.IsMac() {
		return nil
	}

	// Phase 1: Process taps (must be sequential to ensure tap exists before packages)
	var tapPackages []config.BrewSpec
	for _, tap := range group.BrewTapSpec {
		tapName := tap.Name

		// Check if tap is already added
		if err := sh.RunShell(fmt.Sprintf("brew tap | grep -q '^%s$'", tapName)); err == nil {
			onComplete(types.TaskResult{
				Name:   "tap:" + tapName,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already tapped"),
			})
		} else {
			// Add the tap
			tapCmd := fmt.Sprintf("brew tap %s", tapName)
			if tap.URL != nil && *tap.URL != "" {
				tapCmd = fmt.Sprintf("brew tap %s %s", tapName, *tap.URL)
			}

			if err := sh.RunShell(tapCmd); err != nil {
				onComplete(types.TaskResult{
					Name:   "tap:" + tapName,
					Status: types.StatusFailed,
					Error:  err,
				})
				return fmt.Errorf("failed to tap %s: %w", tapName, err)
			}

			onComplete(types.TaskResult{
				Name:   "tap:" + tapName,
				Status: types.StatusSuccess,
				Error:  nil,
			})
		}

		// Collect packages from tap for batch installation
		tapPackages = append(tapPackages, tap.Pkgs...)
	}

	// Phase 2: Collect all packages and group by cask/non-cask
	var formulas []config.BrewSpec
	var casks []config.BrewSpec

	allPackages := append(tapPackages, group.Brew...)
	for _, pkg := range allPackages {
		if pkg.Cask {
			casks = append(casks, pkg)
		} else {
			formulas = append(formulas, pkg)
		}
	}

	// Phase 3: Batch install formulas
	if len(formulas) > 0 {
		if err := p.batchInstallWithVerification(formulas, false, onComplete); err != nil {
			return err
		}
	}

	// Phase 4: Batch install casks
	if len(casks) > 0 {
		if err := p.batchInstallWithVerification(casks, true, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) batchInstallWithVerification(packages []config.BrewSpec, isCask bool, onComplete types.OnTaskComplete) error {
	// Filter out already installed packages
	var toInstall []string
	for _, pkg := range packages {
		brewArgs := ""
		if isCask {
			brewArgs = "--cask "
		}

		// Check if already installed
		if err := sh.RunShell(fmt.Sprintf("brew list %s%s &> /dev/null", brewArgs, pkg.Name)); err == nil {
			onComplete(types.TaskResult{
				Name:   pkg.Name,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
		} else {
			toInstall = append(toInstall, pkg.Name)
		}
	}

	// Nothing to install
	if len(toInstall) == 0 {
		return nil
	}

	// Show progress message
	pkgType := "formulas"
	if isCask {
		pkgType = "casks"
	}
	fmt.Printf("⚙️  Installing %d brew %s: %s\n", len(toInstall), pkgType, strings.Join(toInstall, ", "))

	// Batch install
	brewArgs := ""
	if isCask {
		brewArgs = "--cask "
	}

	installCmd := fmt.Sprintf("HOMEBREW_NO_AUTO_UPDATE=1 brew install %s%s", brewArgs, strings.Join(toInstall, " "))
	if err := sh.RunShell(installCmd); err != nil {
		// On batch install failure, report all as failed
		for _, pkg := range toInstall {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("batch install failed: %w", err),
			})
		}
		return fmt.Errorf("failed to batch install brew packages: %w", err)
	}

	// Post-verification: verify each package individually
	for _, pkg := range toInstall {
		brewArgs := ""
		if isCask {
			brewArgs = "--cask "
		}

		if err := sh.RunShell(fmt.Sprintf("brew list %s%s &> /dev/null", brewArgs, pkg)); err != nil {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("verification failed: %w", err),
			})
			return fmt.Errorf("failed to verify brew package %s: %w", pkg, err)
		}

		onComplete(types.TaskResult{
			Name:   pkg,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}
