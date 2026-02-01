//go:build linux

package apt

import (
	"fmt"
	"os"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/osutil"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct {
	aptUpdated bool
}

func NewProvider() *Provider {
	return &Provider{
		aptUpdated: false,
	}
}

func (p *Provider) ID() string {
	return "apt"
}

func (p *Provider) Priority() int {
	return 10 // System package manager - high priority
}

func (p *Provider) Setup() error {
	// Skip if not Linux
	if !osutil.IsLinux() {
		return nil
	}

	// Ensure apt update is run once
	return p.ensureAptUpdate(false)
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if not Linux
	if !osutil.IsLinux() {
		return nil
	}

	var allPackagesToInstall []string
	ppaSourcesAdded := false

	// Phase 1: Add all PPA sources (no update yet)
	for _, ppa := range group.PPA {
		added, packages, err := p.addPPASource(ppa, onComplete)
		if err != nil {
			return err
		}
		if added {
			ppaSourcesAdded = true
		}
		// Collect packages from this PPA
		for _, pkg := range packages {
			if !p.isAptInstalled(pkg) {
				allPackagesToInstall = append(allPackagesToInstall, pkg)
			}
		}
	}

	// Phase 2: Single apt-update if any PPAs were added
	if ppaSourcesAdded {
		if err := p.ensureAptUpdate(true); err != nil {
			return err
		}
	}

	// Phase 3: Process regular Apt packages
	for _, pkg := range group.Apt {
		if p.isURL(pkg) {
			// URL-based debs still handled individually
			if err := p.installDebFromURL(pkg, onComplete); err != nil {
				return err
			}
		} else if p.isAptInstalled(pkg) {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
		} else {
			allPackagesToInstall = append(allPackagesToInstall, pkg)
		}
	}

	// Phase 4: Batch install all packages at once
	if len(allPackagesToInstall) > 0 {
		if err := p.batchInstallWithVerification(allPackagesToInstall, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) ensureAptUpdate(force bool) error {
	if !p.aptUpdated || force {
		if err := sh.Run("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt index: %w", err)
		}
		p.aptUpdated = true
	}
	return nil
}

func (p *Provider) batchInstallWithVerification(packages []string, onComplete types.OnTaskComplete) error {
	if len(packages) == 0 {
		return nil
	}

	if err := p.ensureAptUpdate(false); err != nil {
		return err
	}

	// Show progress message
	fmt.Printf("⚙️  Installing %d apt packages: %s\n", len(packages), strings.Join(packages, ", "))

	// Execute batch install
	args := append([]string{"apt", "install", "-y"}, packages...)
	err := sh.Run("sudo", args...)

	// Post-verification: Check each package individually
	for _, pkg := range packages {
		if p.isAptInstalled(pkg) {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSuccess,
				Error:  nil,
			})
		} else {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusFailed,
				Error:  fmt.Errorf("package not found after installation"),
			})
		}
	}

	// Return error only if the batch command itself failed critically
	if err != nil {
		return fmt.Errorf("batch apt install encountered errors: %w", err)
	}

	return nil
}

func (p *Provider) installAptPackages(packages []string, onComplete types.OnTaskComplete) error {
	var toInstall []string

	for _, pkg := range packages {
		if p.isURL(pkg) {
			if err := p.installDebFromURL(pkg, onComplete); err != nil {
				return err
			}
			continue
		}

		if p.isAptInstalled(pkg) {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}
		toInstall = append(toInstall, pkg)
	}

	return p.batchInstallWithVerification(toInstall, onComplete)
}

// addPPASource adds the PPA source without installing packages
// Returns: (sourceWasAdded, packagesToInstall, error)
func (p *Provider) addPPASource(ppa config.PPASpec, onComplete types.OnTaskComplete) (bool, []string, error) {
	name := ppa.Name

	// Sanitize name for filename
	sourceName := strings.TrimPrefix(name, "ppa:")
	sourceName = strings.ReplaceAll(sourceName, "/", "-")
	sourceName = strings.ReplaceAll(sourceName, ".", "-")

	// Check if source already exists
	sourcesPath := fmt.Sprintf("/etc/apt/sources.list.d/%s.sources", sourceName)
	if _, err := os.Stat(sourcesPath); err == nil {
		// Source already exists, just return packages
		return false, ppa.Pkgs, nil
	}

	// Default values
	suites := ppa.Suites
	if suites == "" {
		if out, err := sh.RunShellOutput("grep VERSION_CODENAME /etc/os-release | cut -d= -f2"); err == nil {
			suites = strings.TrimSpace(out)
		} else {
			suites = "stable"
		}
	}

	components := ppa.Components
	if components == "" {
		components = "main"
	}

	// Handle GPG Key
	keyringPath := ""
	if ppa.Key != "" || (ppa.KeyServer != "" && ppa.KeyID != "") {
		keyringPath = fmt.Sprintf("/usr/share/keyrings/%s-keyring.gpg", sourceName)
		if _, err := os.Stat(keyringPath); os.IsNotExist(err) {
			if ppa.Key != "" {
				if strings.HasSuffix(ppa.Key, ".asc") || !strings.HasSuffix(ppa.Key, ".gpg") {
					if err := sh.RunShell(fmt.Sprintf("curl -sSL %s | sudo gpg --dearmor --yes --output %s", ppa.Key, keyringPath)); err != nil {
						onComplete(types.TaskResult{
							Name:   "ppa:" + name,
							Status: types.StatusFailed,
							Error:  err,
						})
						return false, nil, fmt.Errorf("failed to download and dearmor GPG key: %w", err)
					}
				} else {
					if err := sh.RunShell(fmt.Sprintf("curl -sS %s | sudo tee %s > /dev/null", ppa.Key, keyringPath)); err != nil {
						onComplete(types.TaskResult{
							Name:   "ppa:" + name,
							Status: types.StatusFailed,
							Error:  err,
						})
						return false, nil, fmt.Errorf("failed to download GPG key: %w", err)
					}
				}
			} else if ppa.KeyServer != "" && ppa.KeyID != "" {
				if err := sh.RunShell(fmt.Sprintf("sudo gpg --no-default-keyring --keyring %s --keyserver %s --recv-keys %s", keyringPath, ppa.KeyServer, ppa.KeyID)); err != nil {
					onComplete(types.TaskResult{
						Name:   "ppa:" + name,
						Status: types.StatusFailed,
						Error:  err,
					})
					return false, nil, fmt.Errorf("failed to download GPG key from server: %w", err)
				}
			}
		}
	}

	// Create .sources file
	finalURI := ppa.URI
	if finalURI == "" && strings.HasPrefix(name, "ppa:") {
		ppaParts := strings.Split(strings.TrimPrefix(name, "ppa:"), "/")
		if len(ppaParts) == 2 {
			finalURI = fmt.Sprintf("https://ppa.launchpadexternal.net/%s/%s/ubuntu", ppaParts[0], ppaParts[1])
		}
	}

	if finalURI == "" {
		err := fmt.Errorf("PPA requires a URI or ppa:user/repo format")
		onComplete(types.TaskResult{
			Name:   "ppa:" + name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return false, nil, err
	}

	// Get architecture
	arch := "amd64"
	if out, err := sh.RunShellOutput("dpkg --print-architecture"); err == nil {
		arch = strings.TrimSpace(out)
	}
	finalURI = strings.ReplaceAll(finalURI, ":arch", arch)

	sourceContent := fmt.Sprintf(`Types: deb
URIs: %s
Suites: %s
Components: %s
Architectures: %s
`, finalURI, suites, components, arch)

	if keyringPath != "" {
		sourceContent += fmt.Sprintf("Signed-By: %s\n", keyringPath)
	}

	// Write to temp file and move
	tmpFile := "/tmp/" + sourceName + ".sources"
	if err := os.WriteFile(tmpFile, []byte(sourceContent), 0644); err != nil {
		onComplete(types.TaskResult{
			Name:   "ppa:" + name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return false, nil, fmt.Errorf("failed to create temp sources file: %w", err)
	}

	if err := sh.Run("sudo", "mv", tmpFile, sourcesPath); err != nil {
		onComplete(types.TaskResult{
			Name:   "ppa:" + name,
			Status: types.StatusFailed,
			Error:  err,
		})
		return false, nil, fmt.Errorf("failed to move sources file: %w", err)
	}

	fmt.Printf("✅ Added PPA source: %s\n", name)

	onComplete(types.TaskResult{
		Name:   "ppa:" + name,
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return true, ppa.Pkgs, nil
}

func (p *Provider) installDebFromURL(url string, onComplete types.OnTaskComplete) error {
	if err := p.ensureAptUpdate(false); err != nil {
		onComplete(types.TaskResult{
			Name:   url,
			Status: types.StatusFailed,
			Error:  err,
		})
		return err
	}

	err := sh.RunShell(fmt.Sprintf("curl -L -o /tmp/package.deb %s && sudo apt install -y /tmp/package.deb && rm /tmp/package.deb", url))
	if err != nil {
		onComplete(types.TaskResult{
			Name:   url,
			Status: types.StatusFailed,
			Error:  err,
		})
		return fmt.Errorf("failed to install deb from url: %w", err)
	}

	onComplete(types.TaskResult{
		Name:   url,
		Status: types.StatusSuccess,
		Error:  nil,
	})

	return nil
}

func (p *Provider) isAptInstalled(pkg string) bool {
	cmd := fmt.Sprintf("dpkg-query -W --showformat='${Status}\\n' %s 2>/dev/null | grep 'install ok installed'", pkg)
	return sh.RunShell(cmd) == nil
}

func (p *Provider) isURL(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}
