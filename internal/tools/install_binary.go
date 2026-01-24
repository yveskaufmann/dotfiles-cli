package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InstallBinary(name, url, version, binaryPath, installPath string, binaries []string) error {
	// Check if already installed
	if len(binaries) > 0 {
		allInstalled := true
		for _, b := range binaries {
			if err := RunShell("type " + b + " > /dev/null 2>&1"); err != nil {
				allInstalled = false
				break
			}
		}
		if allInstalled {
			fmt.Printf("✅ Binaries %v are already installed\n", binaries)
			return nil
		}
	} else if err := RunShell("type " + name + " > /dev/null 2>&1"); err == nil {
		fmt.Printf("✅ %s is already installed (type check passed)\n", name)
		return nil
	}

	if version == "" {
		version = "latest"
	}

	finalURL := strings.ReplaceAll(url, ":version", version)
	fmt.Printf("⚙️  Installing from %s\n", finalURL)

	tmpFile := filepath.Join(os.TempDir(), filepath.Base(finalURL))
	if strings.Contains(tmpFile, "?") {
		tmpFile = filepath.Join(os.TempDir(), name+"_download")
	}
	defer os.Remove(tmpFile)

	fmt.Printf("⬇️  Downloading %s...\n", finalURL)
	if err := RunShell(fmt.Sprintf("curl -L -sS -o %s %s", tmpFile, finalURL)); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	if len(binaries) > 0 {
		return InstallMultipleFromArchive(binaries, tmpFile, installPath)
	}

	return InstallFromArchiveOrBinary(name, tmpFile, binaryPath, installPath)
}

// InstallMultipleFromArchive extracts multiple binaries from an archive and installs them.
func InstallMultipleFromArchive(binaries []string, srcFile, installPath string) error {
	fmt.Printf("📦 Extracting multiple binaries from %s...\n", srcFile)
	extractDir := filepath.Join(os.TempDir(), "multi_extract")
	os.RemoveAll(extractDir)
	os.MkdirAll(extractDir, 0755)
	defer os.RemoveAll(extractDir)

	if strings.HasSuffix(srcFile, ".zip") {
		if err := RunShell(fmt.Sprintf("unzip -o %s -d %s", srcFile, extractDir)); err != nil {
			return fmt.Errorf("failed to extract zip: %w", err)
		}
	} else {
		if err := RunShell(fmt.Sprintf("tar -xzf %s -C %s", srcFile, extractDir)); err != nil {
			return fmt.Errorf("failed to extract tar: %w", err)
		}
	}

	for _, name := range binaries {
		// Try to find the binary by name in the extract dir
		srcBinary := filepath.Join(extractDir, name)
		if _, err := os.Stat(srcBinary); os.IsNotExist(err) {
			// Search recursively
			found := ""
			filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() && info.Name() == name {
					found = path
					return filepath.SkipAll
				}
				return nil
			})
			if found != "" {
				srcBinary = found
			}
		}

		if _, err := os.Stat(srcBinary); os.IsNotExist(err) {
			return fmt.Errorf("binary %s not found in archive", name)
		}

		if err := FinalizeBinaryInstall(srcBinary, installPath, name); err != nil {
			return err
		}
	}

	return nil
}

// InstallFromArchiveOrBinary handles extraction if needed and finalizes the installation.
func InstallFromArchiveOrBinary(name, srcFile, binaryPath, installPath string) error {
	if strings.HasSuffix(srcFile, ".tar.gz") || strings.HasSuffix(srcFile, ".tgz") || strings.HasSuffix(srcFile, ".zip") {
		fmt.Printf("📦 Extracting archive %s...\n", srcFile)
		extractDir := filepath.Join(os.TempDir(), name+"_extract")
		os.RemoveAll(extractDir)
		os.MkdirAll(extractDir, 0755)
		defer os.RemoveAll(extractDir)

		if strings.HasSuffix(srcFile, ".zip") {
			if err := RunShell(fmt.Sprintf("unzip -o %s -d %s", srcFile, extractDir)); err != nil {
				return fmt.Errorf("failed to extract zip: %w", err)
			}
		} else {
			if err := RunShell(fmt.Sprintf("tar -xzf %s -C %s", srcFile, extractDir)); err != nil {
				return fmt.Errorf("failed to extract tar: %w", err)
			}
		}

		// find binary
		srcBinary := filepath.Join(extractDir, binaryPath)
		if binaryPath == "" {
			// Try to find the binary by name in the extract dir
			srcBinary = filepath.Join(extractDir, name)
			if _, err := os.Stat(srcBinary); os.IsNotExist(err) {
				// Search recursively for the first file matching 'name'
				found := ""
				filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
					if err == nil && !info.IsDir() && info.Name() == name {
						found = path
						return filepath.SkipAll
					}
					return nil
				})
				if found != "" {
					srcBinary = found
				}
			}
		}

		// Double check it exists
		if _, err := os.Stat(srcBinary); os.IsNotExist(err) {
			return fmt.Errorf("binary %s not found in archive", name)
		}

		return FinalizeBinaryInstall(srcBinary, installPath, name)
	}

	// Assumed to be a binary already
	return FinalizeBinaryInstall(srcFile, installPath, name)
}

// FinalizeBinaryInstall moves a binary from a temporary path to a target directory.
// If targetDir is empty, it defaults to ~/bin.
func FinalizeBinaryInstall(srcPath, targetDir, name string) error {
	if targetDir == "" {
		targetDir = filepath.Join(os.Getenv("HOME"), "bin")
	}

	destPath := filepath.Join(targetDir, name)

	// Autodetect if sudo is needed (if outside home)
	homeDir := os.Getenv("HOME")
	useSudo := !strings.HasPrefix(targetDir, homeDir)

	mkdirCmd := "mkdir -p"
	moveCmd := "mv"

	if useSudo {
		mkdirCmd = "sudo " + mkdirCmd
		moveCmd = "sudo " + moveCmd
	}

	if err := RunShell(fmt.Sprintf("%s %s", mkdirCmd, targetDir)); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	if err := RunShell(fmt.Sprintf("%s %s %s", moveCmd, srcPath, destPath)); err != nil {
		return fmt.Errorf("failed to move binary to %s: %w", destPath, err)
	}

	if err := RunShell(fmt.Sprintf("sudo chmod +x %s", destPath)); err != nil {
		return fmt.Errorf("failed to make binary %s executable: %w", destPath, err)
	}

	fmt.Printf("✅ Binary %s installed successfully to %s\n", name, destPath)
	return nil
}
