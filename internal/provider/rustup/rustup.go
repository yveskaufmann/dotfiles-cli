package rustup

import (
	"fmt"
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
	return "rustup"
}

func (p *Provider) Priority() int {
	return 100 // Default priority
}

func (p *Provider) HasConfig(group config.DependencyGroup) bool {
	return len(group.RustUp) > 0
}

func (p *Provider) Setup() error {
	if sh.IsBinaryOnPath("rustup") {
		return nil
	}

	fmt.Println("⚙️  rustup not found, installing...")
	if err := sh.RunShell("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y"); err != nil {
		return fmt.Errorf("failed to install rustup: %w", err)
	}

	if err := sh.RunShell("source $HOME/.cargo/env"); err != nil {
		return fmt.Errorf("failed to source rustup environment: %w", err)
	}

	fmt.Println("🚀  rustup installation complete.")
	return nil
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if no rustup packages defined
	if len(group.RustUp) == 0 {
		return nil
	}

	components := make(map[string]bool)
	targets := make(map[string]bool)
	toolchains := make(map[string]bool)
	defaultToolchain := ""

	// Aggregate all specified toolchains, components, and targets
	for _, pkg := range group.RustUp {

		for _, comp := range pkg.Components {
			if ok := components[comp]; ok {
				continue
			}
			components[comp] = true
		}

		for _, target := range pkg.Targets {
			if ok := targets[target]; ok {
				continue
			}
			targets[target] = true
		}

		for _, toolchain := range pkg.Toolchains {
			if ok := toolchains[toolchain]; ok {
				continue
			}
			toolchains[toolchain] = true
		}

		if pkg.DefaultToolchain != "" {
			defaultToolchain = pkg.DefaultToolchain
		}
	}

	if err := p.installComponents(getMapKeys(components), onComplete); err != nil {
		return err
	}

	if err := p.installTargets(getMapKeys(targets), onComplete); err != nil {
		return err
	}

	if err := p.installToolchains(defaultToolchain, getMapKeys(toolchains), onComplete); err != nil {
		return err
	}

	return nil
}

func (p *Provider) installComponents(components []string, onComplete types.OnTaskComplete) error {
	// Install components
	for _, comp := range components {
		taskLabel := fmt.Sprintf("component %s", comp)

		if err := sh.RunShell(fmt.Sprintf("rustup component list --installed | grep %s > /dev/null 2>&1", comp)); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}
		if err := sh.Run("rustup", "component", "add", comp); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("failed to install component %s: %w", comp, err),
			})
			return fmt.Errorf("failed to install component %s: %w", comp, err)
		}

		onComplete(types.TaskResult{
			Name:   taskLabel,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}

func (p *Provider) installTargets(targets []string, onComplete types.OnTaskComplete) error {
	// Install targets
	for _, target := range targets {

		taskLabel := fmt.Sprintf("target %s", target)
		if err := sh.RunShell(fmt.Sprintf("rustup target list --installed | grep %s > /dev/null 2>&1", target)); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}
		if err := sh.Run("rustup", "target", "add", target); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("failed to install target %s: %w", target, err),
			})
			return fmt.Errorf("failed to install target %s: %w", target, err)
		}

		onComplete(types.TaskResult{
			Name:   taskLabel,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}

func (p *Provider) installToolchains(defaultToolchain string, toolchains []string, onComplete types.OnTaskComplete) error {

	fmt.Printf("⚙️  Installing Rust toolchains: %s\n", strings.Join(append([]string{defaultToolchain}, toolchains...), ", "))
	for _, toolchain := range toolchains {

		taskLabel := fmt.Sprintf("toolchain %s", toolchain)
		if err := sh.RunShell(fmt.Sprintf("rustup toolchain list | grep %s > /dev/null 2>&1", toolchain)); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}
		if err := sh.Run("rustup", "toolchain", "install", toolchain); err != nil {
			onComplete(types.TaskResult{
				Name:   taskLabel,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("failed to install toolchain %s: %w", toolchain, err),
			})
			return fmt.Errorf("failed to install toolchain %s: %w", toolchain, err)
		}

		onComplete(types.TaskResult{
			Name:   taskLabel,
			Status: types.StatusSuccess,
			Error:  nil,
		})
	}

	return nil
}

func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
