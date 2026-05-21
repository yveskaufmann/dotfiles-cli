package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yv35.com/dotfiles-cli/internal/util/sh"
)

// DownloadFile downloads a file from URL to a temporary location.
// Returns the temporary file path. Caller is responsible for cleanup.
func DownloadFile(url, name string) (string, error) {
	tmpFile := filepath.Join(os.TempDir(), filepath.Base(url))
	if strings.Contains(tmpFile, "?") {
		tmpFile = filepath.Join(os.TempDir(), name+"_download")
	}

	fmt.Printf("⬇️  Downloading %s...\n", url)
	if err := sh.RunShell(fmt.Sprintf("curl -L -sS -o %s %s", tmpFile, url)); err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}

	return tmpFile, nil
}

// ExtractArchive extracts an archive to a destination directory.
// Automatically detects format based on file extension (.zip, .tar.gz, .tgz).
func ExtractArchive(archivePath, destDir string) error {
	fmt.Printf("📦 Extracting archive %s...\n", archivePath)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}

	if strings.HasSuffix(archivePath, ".zip") {
		if err := sh.RunShell(fmt.Sprintf("unzip -o %s -d %s", archivePath, destDir)); err != nil {
			return fmt.Errorf("failed to extract zip: %w", err)
		}
	} else if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		if err := sh.RunShell(fmt.Sprintf("tar -xzf %s -C %s", archivePath, destDir)); err != nil {
			return fmt.Errorf("failed to extract tar: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}

	return nil
}

// FindBinary searches for a binary by name in a directory tree.
// Returns the full path to the binary, or an error if not found.
func FindBinary(rootDir, binaryName string) (string, error) {
	// Try direct path first
	directPath := filepath.Join(rootDir, binaryName)
	if _, err := os.Stat(directPath); err == nil {
		return directPath, nil
	}

	// Search recursively
	var found string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == binaryName {
			found = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if found == "" {
		return "", fmt.Errorf("binary %s not found in %s", binaryName, rootDir)
	}

	return found, nil
}

// FindBinaries searches for multiple binaries in a directory tree.
// Returns a map of binary name to full path. Returns error if any binary is not found.
func FindBinaries(rootDir string, binaryNames []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, name := range binaryNames {
		path, err := FindBinary(rootDir, name)
		if err != nil {
			return nil, err
		}
		result[name] = path
	}

	return result, nil
}

// InstallBinary moves a binary to target directory with proper permissions.
// Auto-detects sudo requirement based on target path (outside $HOME) and actual
// write permission — handles the case where a previous root-owned file exists
// inside $HOME (e.g. crane installed earlier via sudo).
// If targetDir is empty, defaults to ~/bin.
func InstallBinary(srcPath, targetDir, name string) error {
	if targetDir == "" {
		targetDir = filepath.Join(os.Getenv("HOME"), "bin")
	}

	destPath := filepath.Join(targetDir, name)

	// Primary heuristic: paths outside $HOME or sourced from /var require sudo.
	homeDir := os.Getenv("HOME")
	useSudo := !strings.HasPrefix(targetDir, homeDir) || strings.HasPrefix(srcPath, "/var")

	// Secondary check: even when inside $HOME, an existing file owned by root
	// (or a root-owned directory) will cause "permission denied". Verify we can
	// actually write before committing to a no-sudo run.
	if !useSudo {
		useSudo = !isWritable(targetDir, destPath)
	}

	mkdirCmd := "mkdir -p"
	moveCmd := "mv"

	if useSudo {
		mkdirCmd = "sudo " + mkdirCmd
		moveCmd = "sudo " + moveCmd
	}

	if err := sh.RunShell(fmt.Sprintf("%s %s", mkdirCmd, targetDir)); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	if err := sh.RunShell(fmt.Sprintf("%s %s %s", moveCmd, srcPath, destPath)); err != nil {
		return fmt.Errorf("failed to move binary to %s: %w", destPath, err)
	}

	chmodCmd := "chmod +x"
	if useSudo {
		chmodCmd = "sudo " + chmodCmd
	}

	if err := sh.RunShell(fmt.Sprintf("%s %s", chmodCmd, destPath)); err != nil {
		return fmt.Errorf("failed to make binary %s executable: %w", destPath, err)
	}

	fmt.Printf("✅ Binary %s installed successfully to %s\n", name, destPath)
	return nil
}

// isWritable returns true if the current process can write to destPath, or to
// targetDir (or its nearest existing ancestor) when destPath does not yet exist.
func isWritable(targetDir, destPath string) bool {
	if _, err := os.Stat(destPath); err == nil {
		f, err := os.OpenFile(destPath, os.O_WRONLY, 0)
		if err != nil {
			return false
		}
		f.Close()
		return true
	}
	// Walk up to the nearest existing directory and probe with a temp file.
	dir := targetDir
	for dir != "" && dir != "." && dir != "/" {
		info, err := os.Stat(dir)
		if err != nil {
			dir = filepath.Dir(dir)
			continue
		}
		if !info.IsDir() {
			return false
		}
		tmp := filepath.Join(dir, fmt.Sprintf(".dotfiles_write_test_%d", os.Getpid()))
		f, err := os.Create(tmp)
		if err != nil {
			return false
		}
		f.Close()
		os.Remove(tmp)
		return true
	}
	return false
}

// IsArchive returns true if the file appears to be an archive based on extension.
func IsArchive(filePath string) bool {
	return strings.HasSuffix(filePath, ".tar.gz") ||
		strings.HasSuffix(filePath, ".tgz") ||
		strings.HasSuffix(filePath, ".zip")
}
