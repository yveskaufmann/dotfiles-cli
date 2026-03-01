package jetbrains

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"yv35.com/dotfiles/internal/util/archive"
)

// installIDE downloads and installs the IDE, choosing the right strategy
// based on the download URL (.dmg on macOS, .tar.gz on Linux).
func (p *Provider) installIDE(candidate installCandidate) error {
	if strings.HasSuffix(candidate.Release.DownloadURL, ".dmg") {
		return p.installFromDmg(candidate)
	}
	return p.installFromTarGz(candidate)
}

// installFromTarGz handles Linux: extract tarball → move to /opt/<ide-dir>.
func (p *Provider) installFromTarGz(candidate installCandidate) error {
	spec := candidate.Spec
	release := candidate.Release
	targetDir := candidate.Dir

	displayName := spec.IDE
	if spec.Name != "" {
		displayName = spec.Name
	}

	tmpFile, err := archive.DownloadFile(release.DownloadURL, displayName)
	if err != nil {
		return fmt.Errorf("failed to download %s %s: %w", displayName, release.Version, err)
	}
	defer os.Remove(tmpFile)

	extractDir := fmt.Sprintf("/tmp/jetbrains_%s_%d", spec.IDE, os.Getpid())
	os.RemoveAll(extractDir)
	defer os.RemoveAll(extractDir)

	if err := archive.ExtractArchive(tmpFile, extractDir); err != nil {
		return fmt.Errorf("failed to extract %s: %w", displayName, err)
	}

	// The tarball always contains exactly one top-level directory
	// (e.g. idea-IU-2024.3.4/). Find it.
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("failed to read extract dir: %w", err)
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return fmt.Errorf("unexpected archive layout for %s: expected exactly one top-level directory", displayName)
	}
	extractedDir := filepath.Join(extractDir, entries[0].Name())

	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return fmt.Errorf("failed to create parent directory %s: %w", filepath.Dir(targetDir), err)
	}
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove existing installation at %s: %w", targetDir, err)
	}
	if err := os.Rename(extractedDir, targetDir); err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", extractedDir, targetDir, err)
	}

	return nil
}

// installFromDmg handles macOS: mount .dmg → copy .app to /Applications/ → unmount.
func (p *Provider) installFromDmg(candidate installCandidate) error {
	spec := candidate.Spec
	release := candidate.Release
	targetDir := candidate.Dir // e.g. /Applications/IntelliJ IDEA.app

	displayName := spec.IDE
	if spec.Name != "" {
		displayName = spec.Name
	}

	tmpFile, err := archive.DownloadFile(release.DownloadURL, displayName)
	if err != nil {
		return fmt.Errorf("failed to download %s %s: %w", displayName, release.Version, err)
	}
	defer os.Remove(tmpFile)

	// Mount the DMG quietly at a deterministic mount point.
	mountPoint := fmt.Sprintf("/Volumes/JetBrains_%s_%d", spec.IDE, os.Getpid())
	attachCmd := exec.Command("hdiutil", "attach", "-nobrowse", "-quiet", "-mountpoint", mountPoint, tmpFile)
	if out, err := attachCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to mount %s: %w\n%s", displayName, err, out)
	}
	defer func() {
		exec.Command("hdiutil", "detach", "-quiet", mountPoint).Run() //nolint:errcheck
	}()

	// Find the .app bundle inside the mounted volume.
	appSrc, err := findAppBundle(mountPoint)
	if err != nil {
		return fmt.Errorf("failed to find .app in DMG for %s: %w", displayName, err)
	}

	// Remove any previous installation, then copy the bundle into /Applications/.
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove existing installation at %s: %w", targetDir, err)
	}
	copyCmd := exec.Command("cp", "-R", appSrc, targetDir)
	if out, err := copyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to copy %s to %s: %w\n%s", appSrc, targetDir, err, out)
	}

	return nil
}

// findAppBundle locates the first .app directory at the top level of dir.
func findAppBundle(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasSuffix(e.Name(), ".app") {
			return filepath.Join(dir, e.Name()), nil
		}
	}
	return "", fmt.Errorf("no .app bundle found in %s", dir)
}
