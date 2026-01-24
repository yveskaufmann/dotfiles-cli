//go:build linux

package tools

import (
	"fmt"
	"os"
	"strings"
)

func InstallAptPackages(packages []string) error {
	for _, pkg := range packages {
		if err := Install(pkg); err != nil {
			return err
		}
	}
	return nil
}

func Install(pkg string) error {
	// check if it is a URL or a package name
	if isURL(pkg) {
		return installDebFromURL(pkg)
	}

	if isAptInstalled(pkg) {
		fmt.Printf("✅ Apt package %s is already installed\n", pkg)
		return nil
	}

	fmt.Printf("⚙️  Installing apt package: %s\n", pkg)
	err := Run("sudo", "apt", "install", "-y", pkg)
	if err != nil {
		return fmt.Errorf("❌ Failed to install %s via apt: %w", pkg, err)
	}
	fmt.Printf("✅ Apt package %s installed successfully\n", pkg)
	return nil
}

func isAptInstalled(pkg string) bool {
	// dpkg-query -W --showformat='${Status}\n' $pkg 2>/dev/null | grep 'install ok installed'
	cmd := fmt.Sprintf("dpkg-query -W --showformat='${Status}\\n' %s 2>/dev/null | grep 'install ok installed'", pkg)
	return RunShell(cmd) == nil
}

func isURL(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}

func InstallPPA(name, keyURL, keyServer, keyID, uri, suites, components string, pkgs []string) error {
	// 1. Sanitize name for filename
	sourceName := strings.TrimPrefix(name, "ppa:")
	sourceName = strings.ReplaceAll(sourceName, "/", "-")
	sourceName = strings.ReplaceAll(sourceName, ".", "-")

	// 2. Default values
	if suites == "" {
		// Try to get Ubuntu codename if suites is empty
		if out, err := RunShellOutput("lsb_release -cs"); err == nil {
			suites = strings.TrimSpace(out)
		} else {
			suites = "stable"
		}
	}
	if components == "" {
		components = "main"
	}

	// 3. Check if packages are already installed
	allInstalled := true
	for _, pkg := range pkgs {
		if !isAptInstalled(pkg) {
			allInstalled = false
			break
		}
	}
	if allInstalled && len(pkgs) > 0 {
		fmt.Printf("✅ PPA packages %v are already installed\n", pkgs)
		return nil
	}

	// 4. Handle GPG Key
	keyringPath := ""
	if keyURL != "" || (keyServer != "" && keyID != "") {
		keyringPath = fmt.Sprintf("/usr/share/keyrings/%s-keyring.gpg", sourceName)
		if _, err := os.Stat(keyringPath); os.IsNotExist(err) {
			if keyURL != "" {
				fmt.Printf("⚙️  Downloading GPG key: %s\n", keyURL)
				if strings.HasSuffix(keyURL, ".asc") {
					if err := RunShell(fmt.Sprintf("curl -sS %s | sudo gpg --dearmor --yes --output %s", keyURL, keyringPath)); err != nil {
						return fmt.Errorf("failed to download and dearmor GPG key %s: %w", keyURL, err)
					}
				} else {
					if err := RunShell(fmt.Sprintf("curl -sS %s | sudo tee %s > /dev/null", keyURL, keyringPath)); err != nil {
						return fmt.Errorf("failed to download GPG key %s: %w", keyURL, err)
					}
				}
			} else if keyServer != "" && keyID != "" {
				fmt.Printf("⚙️  Downloading GPG key from server: %s (ID: %s)\n", keyServer, keyID)
				if err := RunShell(fmt.Sprintf("sudo gpg --no-default-keyring --keyring %s --keyserver %s --recv-keys %s", keyringPath, keyServer, keyID)); err != nil {
					return fmt.Errorf("failed to download GPG key from server %s: %w", keyServer, err)
				}
			}
		}
	}

	// 5. Handle .sources file (DEB822 format)
	sourcesPath := fmt.Sprintf("/etc/apt/sources.list.d/%s.sources", sourceName)
	if _, err := os.Stat(sourcesPath); os.IsNotExist(err) {
		fmt.Printf("⚙️  Adding PPA source: %s\n", sourcesPath)

		finalURI := uri
		if finalURI == "" && strings.HasPrefix(name, "ppa:") {
			// Construct PPA URL if it's a launchpad PPA
			// ppa:user/repo -> http://ppa.launchpad.net/user/repo/ubuntu
			ppaParts := strings.Split(strings.TrimPrefix(name, "ppa:"), "/")
			if len(ppaParts) == 2 {
				finalURI = fmt.Sprintf("https://ppa.launchpadexternal.net/%s/%s/ubuntu", ppaParts[0], ppaParts[1])
			}
		}

		if finalURI == "" {
			return fmt.Errorf("PPA %s requires a URI or have to be a ppa:user/repo format", name)
		}

		// Get architecture
		arch := "amd64"
		if out, err := RunShellOutput("dpkg --print-architecture"); err == nil {
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

		// Write to temp file and move to /etc/apt/sources.list.d/
		tmpFile := "/tmp/" + sourceName + ".sources"
		if err := os.WriteFile(tmpFile, []byte(sourceContent), 0644); err != nil {
			return fmt.Errorf("failed to create temp sources file: %w", err)
		}

		if err := Run("sudo", "mv", tmpFile, sourcesPath); err != nil {
			return fmt.Errorf("failed to move sources file to /etc/apt/sources.list.d/: %w", err)
		}

		if err := Run("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt after adding ppa %s: %w", name, err)
		}
	}

	// 6. Install packages
	return InstallAptPackages(pkgs)
}

func installDebFromURL(url string) error {
	fmt.Printf("⚙️  Installing deb package from URL: %s\n", url)
	// Simple implementation for now
	err := RunShell(fmt.Sprintf("curl -L -o /tmp/package.deb %s && sudo apt install -y /tmp/package.deb && rm /tmp/package.deb", url))
	if err != nil {
		return fmt.Errorf("failed to install deb from url %s: %w", url, err)
	}
	return nil
}
