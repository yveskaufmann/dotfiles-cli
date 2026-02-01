//go:build linux

package tool

import (
	"fmt"
	"os"
	"strings"

	"yv35.com/dotfiles/internal/util/sh"
)

var aptUpdated = false

func ensureAptUpdate(force bool) error {
	if !aptUpdated || force {
		fmt.Println("⚙️  Updating apt index...")
		if err := sh.Run("sudo", "apt", "update"); err != nil {
			return fmt.Errorf("failed to update apt index: %w", err)
		}
		aptUpdated = true
	}
	return nil
}

func InstallAptPackages(packages []string) error {
	var toInstall []string
	for _, pkg := range packages {
		if isURL(pkg) {
			if err := installDebFromURL(pkg); err != nil {
				return err
			}
			continue
		}

		if isAptInstalled(pkg) {
			fmt.Printf("✅ Apt package %s is already installed\n", pkg)
			continue
		}
		toInstall = append(toInstall, pkg)
	}

	if len(toInstall) == 0 {
		return nil
	}

	if err := ensureAptUpdate(false); err != nil {
		return err
	}

	fmt.Printf("⚙️  Installing apt packages: %s\n", strings.Join(toInstall, " "))
	args := append([]string{"apt", "install", "-y"}, toInstall...)
	err := sh.Run("sudo", args...)
	if err != nil {
		return fmt.Errorf("❌ Failed to install packages via apt: %w", err)
	}
	fmt.Printf("✅ Apt packages installed successfully: %s\n", strings.Join(toInstall, ", "))
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

	if err := ensureAptUpdate(false); err != nil {
		return err
	}

	fmt.Printf("⚙️  Installing apt package: %s\n", pkg)
	err := sh.Run("sudo", "apt", "install", "-y", pkg)
	if err != nil {
		return fmt.Errorf("❌ Failed to install %s via apt: %w", pkg, err)
	}
	fmt.Printf("✅ Apt package %s installed successfully\n", pkg)
	return nil
}

func isAptInstalled(pkg string) bool {
	// dpkg-query -W --showformat='${Status}\n' $pkg 2>/dev/null | grep 'install ok installed'
	cmd := fmt.Sprintf("dpkg-query -W --showformat='${Status}\\n' %s 2>/dev/null | grep 'install ok installed'", pkg)
	return sh.RunShell(cmd) == nil
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
		if out, err := sh.RunShellOutput("grep VERSION_CODENAME /etc/os-release | cut -d= -f2"); err == nil {
			suites = strings.TrimSpace(out)
		} else {
			suites = "stable"
			fmt.Printf("⚠️  Unable to determine Ubuntu codename, defaulting suites to 'stable': %v\n", err)
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
				if strings.HasSuffix(keyURL, ".asc") || !strings.HasSuffix(keyURL, ".gpg") {
					if err := sh.RunShell(fmt.Sprintf("curl -sSL %s | sudo gpg --dearmor --yes --output %s", keyURL, keyringPath)); err != nil {
						return fmt.Errorf("failed to download and dearmor GPG key %s: %w", keyURL, err)
					}
				} else {
					if err := sh.RunShell(fmt.Sprintf("curl -sS %s | sudo tee %s > /dev/null", keyURL, keyringPath)); err != nil {
						return fmt.Errorf("failed to download GPG key %s: %w", keyURL, err)
					}
				}
			} else if keyServer != "" && keyID != "" {
				fmt.Printf("⚙️  Downloading GPG key from server: %s (ID: %s)\n", keyServer, keyID)
				if err := sh.RunShell(fmt.Sprintf("sudo gpg --no-default-keyring --keyring %s --keyserver %s --recv-keys %s", keyringPath, keyServer, keyID)); err != nil {
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

		fmt.Printf("%s\n", sourceContent)

		// Write to temp file and move to /etc/apt/sources.list.d/
		tmpFile := "/tmp/" + sourceName + ".sources"
		if err := os.WriteFile(tmpFile, []byte(sourceContent), 0644); err != nil {
			return fmt.Errorf("failed to create temp sources file: %w", err)
		}

		if err := sh.Run("sudo", "mv", tmpFile, sourcesPath); err != nil {
			return fmt.Errorf("failed to move sources file to /etc/apt/sources.list.d/: %w", err)
		}

		if err := ensureAptUpdate(true); err != nil {
			return fmt.Errorf("failed to update apt after adding ppa %s: %w", name, err)
		}
	}

	// 6. Install packages
	return InstallAptPackages(pkgs)
}

func installDebFromURL(url string) error {
	fmt.Printf("⚙️  Installing deb package from URL: %s\n", url)

	if err := ensureAptUpdate(false); err != nil {
		return err
	}

	// Simple implementation for now
	err := sh.RunShell(fmt.Sprintf("curl -L -o /tmp/package.deb %s && sudo apt install -y /tmp/package.deb && rm /tmp/package.deb", url))
	if err != nil {
		return fmt.Errorf("failed to install deb from url %s: %w", url, err)
	}
	return nil
}
